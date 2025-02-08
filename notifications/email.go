package notifications

import (
	"ecommerce/models"
	"fmt"

	"github.com/wneessen/go-mail"
)

// EmailConfig holds the configuration for sending emails.
type EmailConfig struct {
	SMTPHost string
	Port     int
	Username string
	Password string
	From     string
}

func NewEmailConfig(port int, smtpHost, userName, password, from string) *EmailConfig {
	return &EmailConfig{
		SMTPHost: smtpHost,
		Port:     port,
		Username: userName,
		Password: password,
		From:     from,
	}
}

type EmaiMetadata struct {
	To   string
	Body string
}

type NotificationMetadata struct {
	UserID      string
	OrderID     string
	OrderStatus string
	UserEmail   string
}

// SendEmail sends an email notification to the specified recipient.
func (e *EmailConfig) sendEmail(emailMetadata *EmaiMetadata) error {

	m := mail.NewMsg()
	m.From(e.From)
	m.To(emailMetadata.To)
	m.Subject("Order Status")
	m.SetBodyString("text/plain", emailMetadata.Body)

	// Configure and send
	c, err := mail.NewClient(e.SMTPHost, mail.WithPort(e.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(e.Username), mail.WithPassword(e.Password))
	if err != nil {
		fmt.Println("Error creating mail client:", err)
		return err
	}
	err = c.DialAndSend(m)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Email sent successfully")
	return nil
}

func (e *EmailConfig) NotifyOrderStatus(OrderDetails *models.OrderDetails, status string) error {
	body := fmt.Sprintf(`
Hi %s,

Your order has been placed successfully! 
Thank you for shopping with us.

Your order details:
	Order ID: %s
	Order Status: %s
	Total Price: Rs %.2f

We will notify you once your order is processed and shipped.
If you have any queries, feel free to contact us.

Best regards,
Your E-Commerce Team
	`, OrderDetails.Username, OrderDetails.ID, status, OrderDetails.TotalPrice)

	emailMetadata := EmaiMetadata{
		To:   OrderDetails.Email,
		Body: body,
	}

	err := e.sendEmail(&emailMetadata)
	return err
}

func (e *EmailConfig) NotifyUserCreated(userInfo *models.User) error {
	body := fmt.Sprintf(`
Hi %s,

Welcome to Ecommerce services!

We’re thrilled to have you on board. Here’s what you need to know:
Your account has been successfully created.
You can log in anytime using your registered email address: %s

If you have any questions or need assistance, feel free to contact us.

We’re excited to have you with us and look forward to helping you make the most of our services.

Best regards,  
Ecommerce Team  
		`, userInfo.FirstName, userInfo.Email)

	emaiMetadata := EmaiMetadata{
		To:   userInfo.Email,
		Body: body,
	}

	err := e.sendEmail(&emaiMetadata)
	return err
}

// func (e *EmaiMetadata) SendNotification(notificationMetadata *NotificationMetadata) error{}
