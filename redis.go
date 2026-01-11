package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
    defaultHistoryLimit = 100
)

func channelMessagesKey(channel string) string {
    return "channel:" + channel + ":messages"
}

func channelPubsubTopic(channel string) string {
    return "channel:" + channel
}

// sendMessage adds the message to a ZSET (scored by timestamp) and publishes
// it on the channel's pubsub topic for realtime delivery.
func sendMessage(rdb *redis.Client, channel string, msg Message) error {
    ctx := context.Background()

    // Ensure timestamp is set
    if msg.Timestamp == 0 {
        msg.Timestamp = time.Now().Unix()
    }

    // Serialize message
    payload, err := json.Marshal(msg)
    if err != nil {
        return err
    }

    // Store in ZSET
    if err := rdb.ZAdd(ctx, channelMessagesKey(channel), &redis.Z{Score: float64(msg.Timestamp), Member: payload}).Err(); err != nil {
        return err
    }

    // Publish for subscribers
    if err := rdb.Publish(ctx, channelPubsubTopic(channel), payload).Err(); err != nil {
        return err
    }

    return nil
}

// fetchRecentMessages returns up to defaultHistoryLimit most recent messages
// for the given channel, ordered by ascending timestamp.
func fetchRecentMessages(rdb *redis.Client, channel string) ([]Message, error) {
    ctx := context.Background()

    // Fetch latest N entries from ZSET
    vals, err := rdb.ZRevRange(ctx, channelMessagesKey(channel), 0, defaultHistoryLimit-1).Result()
    if err != nil {
        return nil, err
    }

    // Reverse to ascending chronological order
    n := len(vals)
    msgs := make([]Message, 0, n)
    for i := n - 1; i >= 0; i-- {
        var m Message
        if err := json.Unmarshal([]byte(vals[i]), &m); err == nil {
            msgs = append(msgs, m)
        }
    }
    return msgs, nil
}

// subscribeChannel starts a Redis pubsub subscription for the channel and
// forwards decoded messages to the provided sink channel. It runs until the
// provided context is canceled.
func subscribeChannel(ctx context.Context, rdb *redis.Client, channel string, sink chan<- Message) error {
    sub := rdb.Subscribe(ctx, channelPubsubTopic(channel))
    ch := sub.Channel()

    go func() {
        defer sub.Close()
        for {
            select {
            case <-ctx.Done():
                return
            case m, ok := <-ch:
                if !ok {
                    return
                }
                var msg Message
                if err := json.Unmarshal([]byte(m.Payload), &msg); err == nil {
                    sink <- msg
                }
            }
        }
    }()

    return nil
}

// tea helper: wait for one message from a channel and deliver it as a tea.Msg.
// Defined here to keep Redis-related plumbing in one place.
type incomingMsg struct{ Message Message }
