// package service
//
// import (
//   "github.com/HarshithRajesh/app-chat/internal/repository"
//   "context"
//   "github.com/redis/go-redis/v9"
//   "fmt"
//   "time"
// )
//
// func StartMessageConsumer(ctx context.Context,redisClient *redis.Client,streamName string,groupName string,
//                           consumerName string,count int64,block time.Duration){
//   for {
//     res,err := repository.ReadMessagesFromGroup(ctx,redisClient,streamName,groupName,consumerName,count,block)
//
//     if err != nil{
//       if err == redis.Nil{
//         fmt.Printf("No new messages from the group")
//         continue 
//       }
//       fmt.Printf("Error reading from group: %v\n", err)
//       continue
//     }
//     var idsToAck []string
//     if len(res)>0{
//       for _,stream:= range res{
//          for _,msg := range stream.Messages{
//            fmt.Printf("Message Id : %s\n",msg.ID)
//            fmt.Println("Values:")
//            for field,value := range msg.Values{
//               fmt.Printf(" %s: %v\n",field,value)
//             }
//             fmt.Println()
//             idsToAck=append(idsToAck,msg.ID)
//
//
//             }
//          }
//       }
//         if len(idsToAck)>0{
//               ackCount,err := repository.AcknowledgeMessages(ctx,redisClient,streamName,groupName,idsToAck)
//               if err != nil{
//                 fmt.Printf("Error acknowledging %d messages from group %s: %v\n", len(idsToAck), groupName, err)
//               }else{
//                 fmt.Printf("Successfully acknowledged %d messages from group %s \n",ackCount,groupName)
//               }
//   }
//   }
// }
package service // Make sure this matches your package name

import (
	"context"
	"fmt"
	"time" // Needed for time.Duration
	"log" // Good practice for logging fatal errors

	"github.com/redis/go-redis/v9" // Your Redis client library
	"github.com/HarshithRajesh/app-chat/internal/repository" // Your repository package
)

func StartMessageConsumer(ctx context.Context, redisClient *redis.Client, streamName string, groupName string,
	consumerName string, count int64, block time.Duration) {
	log.Printf("Starting consumer goroutine for group %s, consumer %s on stream %s", groupName, consumerName, streamName) // Log consumer start

	for { // Outer continuous loop
		// Call repository to read messages from the group (and block)
		res, err := repository.ReadMessagesFromGroup(ctx, redisClient, streamName, groupName, consumerName, count, block)

		// --- Check for Context Cancellation (for graceful shutdown) ---
		select {
		case <-ctx.Done():
			log.Printf("Consumer for group %s, consumer %s shutting down due to context cancellation", groupName, consumerName)
			return // Exit the goroutine/function cleanly
		default:
			// Continue if context is not done
		}
        // --- End Context Check ---


		if err != nil {
			if err == redis.Nil {
				// Timeout, no new messages during the block duration. Just continue the loop.
				// fmt.Printf("No new messages from the group (timeout)\n") // Optional: Log timeout if desired
				continue // Skip to the next loop iteration to read again
			}
			// Other error reading from Redis - Log it and continue the loop after a pause.
			log.Printf("Error reading from group %s, consumer %s on stream %s: %v. Retrying in 5 seconds.", groupName, consumerName, streamName, err)
			time.Sleep(5 * time.Second) // Pause before retrying on non-timeout error
			continue // Skip to the next loop iteration to read again
		}

		var idsToAck []string // Slice to collect IDs for messages processed in THIS batch

		// --- Message Processing ---
        // Loop through the received stream results (usually one per stream name)
		for _, stream := range res {
            // Loop through the actual messages within the stream result
			for _, msg := range stream.Messages {
				// --- Process each individual message here ---
				// In a real app: validate data, save to DB, send notification, etc.
				// For now, we're just printing:
				fmt.Printf("Consumer %s received message %s:\n", consumerName, msg.ID)
				fmt.Println("  Values:")
				for field, value := range msg.Values {
					fmt.Printf("    %s: %v\n", field, value)
				}
				fmt.Println() // Newline after processing message

				// --- Collect the ID of the processed message for later acknowledgment ---
				idsToAck = append(idsToAck, msg.ID)
			}
		}
		// --- End Message Processing ---


		// --- ACKNOWLEDGE MESSAGES HERE, AFTER PROCESSING THE ENTIRE BATCH ---
        // This code runs AFTER all messages from the `res` batch have been processed in the loops above.
		if len(idsToAck) > 0 { // Check if we actually processed any messages in this batch
			// Call the repository function to acknowledge the collected message IDs
			ackCount, ackErr := repository.AcknowledgeMessages(ctx, redisClient, streamName, groupName, idsToAck) // Correct: Pass the slice idsToAck!
			if ackErr != nil {
				// Log the acknowledgment error. This is important to monitor.
				// The consumer typically should NOT panic here, as it might leave other messages unacknowledged.
				log.Printf("Error acknowledging %d messages from group %s on stream %s: %v", len(idsToAck), groupName, streamName, ackErr)
				// Messages not acknowledged will remain in PEL and might be redelivered later.
			} else {
                // Optional: Log successful acknowledgment
                log.Printf("Successfully acknowledged %d messages from group %s on stream %s", ackCount, groupName, streamName)
            }
		}
        // --- End Acknowledgment ---


		// After processing the batch and acknowledging, the for {} loop continues to the next iteration,
        // calling ReadMessagesFromGroup again to fetch the next batch (and block if needed).
	} // End for {} loop
}

// Consider adding a logging library (like logrus or zap) instead of fmt.Printf/log.Printf for better production logging.
// Also added a basic context check within the loop for graceful shutdown when ctx is cancelled.
