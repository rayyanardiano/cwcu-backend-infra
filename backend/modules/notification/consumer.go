package notification

import (
	"encoding/json"
	"log"

	"cwcu/backend/config"
	"cwcu/backend/pkg/lavinmq"
)

// StartConsumers starts all LavinMQ consumers in background goroutines
func StartConsumers(mq *lavinmq.Client, cfg *config.Config) {
	go consumeNotifications(mq, "notification.email", handleEmailNotification)
	go consumeNotifications(mq, "notification.push", handlePushNotification)
	go consumeNotifications(mq, "grade.updated", handleGradeUpdated)
	go consumeNotifications(mq, "krs.approved", handleKRSApproved)
	go consumeNotifications(mq, "assignment.submitted", handleAssignmentSubmitted)

	log.Println("✅ Notification consumers started")
}

func consumeNotifications(mq *lavinmq.Client, queue string, handler func([]byte)) {
	deliveries, err := mq.Consume(queue, queue+"-consumer")
	if err != nil {
		log.Printf("❌ Failed to consume queue %s: %v", queue, err)
		return
	}

	for d := range deliveries {
		handler(d.Body)
		d.Ack(false) // manual ack
	}
}

func handleEmailNotification(body []byte) {
	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("❌ Email notification parse error: %v", err)
		return
	}
	// TODO: integrate with email service (SendGrid, SMTP, etc.)
	log.Printf("📧 Email notification: %v", payload)
}

func handlePushNotification(body []byte) {
	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("❌ Push notification parse error: %v", err)
		return
	}
	// TODO: integrate with WebSocket or FCM
	log.Printf("🔔 Push notification: event=%s", payload["event"])
}

func handleGradeUpdated(body []byte) {
	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		return
	}
	log.Printf("📊 Grade updated for student=%s class=%s", payload["student_id"], payload["class_id"])
}

func handleKRSApproved(body []byte) {
	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		return
	}
	log.Printf("✅ KRS approved: krs_id=%s", payload["krs_id"])
}

func handleAssignmentSubmitted(body []byte) {
	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		return
	}
	log.Printf("📝 Assignment submitted: assignment_id=%s by student=%s", payload["assignment_id"], payload["student_id"])
}
