import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";

export default function Chat() {
  const [selectedContact, setSelectedContact] = useState(null);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isSending, setIsSending] = useState(false);
  const [error, setError] = useState("");
  
  const navigate = useNavigate();
  const messagesEndRef = useRef(null);
  const pollingInterval = useRef(null);

  // Auto-scroll to bottom when new messages arrive
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    // Get the selected contact from localStorage
    const contactData = localStorage.getItem('selectedContact');
    if (contactData) {
      const contact = JSON.parse(contactData);
      setSelectedContact(contact);
      console.log("Chat with:", contact);
      
      // Fetch existing messages for this contact
      fetchMessages(contact.id);
      
      // Start polling for new messages every 3 seconds
      startPolling(contact.id);
    } else {
      // No contact selected, redirect to contacts page
      navigate('/contact');
    }

    // Cleanup polling when component unmounts
    return () => {
      if (pollingInterval.current) {
        clearInterval(pollingInterval.current);
      }
    };
  }, [navigate]);

  const startPolling = (contactId) => {
    // Clear any existing polling
    if (pollingInterval.current) {
      clearInterval(pollingInterval.current);
    }

    // Poll for new messages every 3 seconds
    pollingInterval.current = setInterval(() => {
      console.log("Polling for new messages...");
      fetchMessages(contactId, true); // true = silent refresh (no loading state)
    }, 3000);
  };

  const parseTimestamp = (timestamp) => {
    // Handle different timestamp formats
    if (!timestamp) return new Date();
    
    // If it's already a Date object
    if (timestamp instanceof Date) return timestamp;
    
    // If it's a string, parse it
    if (typeof timestamp === 'string') {
      // Handle ISO format like "2025-04-10T23:06:52.391458Z"
      return new Date(timestamp);
    }
    
    // Fallback
    return new Date(timestamp);
  };

  const fetchMessages = async (contactId, silentRefresh = false) => {
    if (!silentRefresh) {
      setIsLoading(true);
    }
    setError("");

    try {
      // Get current user ID from localStorage
      const userData = JSON.parse(localStorage.getItem('user') || '{}');
      const userId = userData.id;

      if (!userId) {
        setError("User not found. Please log in again.");
        if (!silentRefresh) setIsLoading(false);
        return;
      }

      console.log("Fetching message history between user:", userId, "and contact:", contactId);

      const response = await fetch(`http://localhost:8080/user/message/history?user1=${userId}&user2=${contactId}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      });

      console.log("Fetch messages response status:", response.status);

      if (response.ok) {
        const data = await response.json();
        console.log("Message history response:", data);
        
        // Process the messages from backend
        if (data.data && Array.isArray(data.data)) {
          const processedMessages = data.data.map(msg => {
            const parsedTimestamp = parseTimestamp(msg.CreatedAt || msg.created_at || msg.timestamp);
            
            return {
              id: msg.Id || msg.id || `${msg.SenderId}-${msg.ReceiverId}-${Date.now()}`, // Ensure unique ID
              content: msg.Content || msg.content || '',
              sender_id: parseInt(msg.SenderId || msg.sender_id),
              receiver_id: parseInt(msg.ReceiverId || msg.receiver_id),
              timestamp: parsedTimestamp,
              isMe: parseInt(msg.SenderId || msg.sender_id) === parseInt(userId),
              // Add raw timestamp for comparison
              rawTimestamp: msg.CreatedAt || msg.created_at || msg.timestamp
            };
          });
          
          // Sort messages by timestamp (oldest first)
          processedMessages.sort((a, b) => a.timestamp - b.timestamp);
          
          console.log("Processed messages:", processedMessages);
          
          // Update messages - use a more reliable comparison
          setMessages(prevMessages => {
            // Compare by message count and latest message timestamp
            if (prevMessages.length !== processedMessages.length) {
              console.log("Message count changed:", prevMessages.length, "->", processedMessages.length);
              return processedMessages;
            }
            
            // Check if the latest message is different
            if (processedMessages.length > 0 && prevMessages.length > 0) {
              const latestNew = processedMessages[processedMessages.length - 1];
              const latestOld = prevMessages[prevMessages.length - 1];
              
              if (latestNew.rawTimestamp !== latestOld.rawTimestamp || 
                  latestNew.content !== latestOld.content) {
                console.log("New message detected!");
                return processedMessages;
              }
            }
            
            return prevMessages;
          });
          
          if (!silentRefresh) {
            console.log("Loaded", processedMessages.length, "messages");
          }
        } else {
          console.log("No messages in data");
          setMessages([]);
        }
      } else {
        const errorText = await response.text();
        console.log("Fetch messages error:", errorText);
        
        if (response.status !== 404 && !silentRefresh) {
          setError("Failed to load message history");
        } else {
          setMessages([]);
        }
      }
    } catch (error) {
      console.error("Error fetching messages:", error);
      if (!silentRefresh) {
        setError("Failed to load messages. Please try again.");
      }
    }

    if (!silentRefresh) {
      setIsLoading(false);
    }
  };

  const sendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim()) return;
    
    setIsSending(true);
    setError("");

    try {
      // Get current user ID from localStorage
      const userData = JSON.parse(localStorage.getItem('user') || '{}');
      const senderId = userData.id;

      if (!senderId) {
        setError("User not found. Please log in again.");
        setIsSending(false);
        return;
      }

      console.log("Sending message:", newMessage, "from:", senderId, "to:", selectedContact.id);

      const messagePayload = {
        sender_id: parseInt(senderId),
        receiver_id: parseInt(selectedContact.id),
        content: newMessage.trim()
      };

      const response = await fetch('http://localhost:8080/user/message', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(messagePayload)
      });

      console.log("Send message response status:", response.status);

      if (response.ok) {
        console.log("Message sent successfully");
        
        // Clear the input immediately
        setNewMessage("");
        
        // Wait a bit then refresh to ensure backend has processed the message
        setTimeout(() => {
          fetchMessages(selectedContact.id, true);
        }, 500);
        
      } else {
        const errorText = await response.text();
        console.log("Send message error:", errorText);
        setError(errorText || "Failed to send message");
      }
    } catch (error) {
      console.error("Error sending message:", error);
      setError("Something went wrong. Please try again.");
    }

    setIsSending(false);
  };

  const goBackToContacts = () => {
    // Clear polling when leaving chat
    if (pollingInterval.current) {
      clearInterval(pollingInterval.current);
    }
    navigate('/contact');
  };

  const formatTime = (timestamp) => {
    try {
      const date = new Date(timestamp);
      // Check if date is valid
      if (isNaN(date.getTime())) {
        return "Invalid time";
      }
      
      const now = new Date();
      const diff = now - date;
      const seconds = Math.floor(diff / 1000);
      const minutes = Math.floor(seconds / 60);
      const hours = Math.floor(minutes / 60);
      const days = Math.floor(hours / 24);
      
      if (days > 0) {
        return date.toLocaleDateString();
      } else if (hours > 0) {
        return `${hours}h ago`;
      } else if (minutes > 0) {
        return `${minutes}m ago`;
      } else {
        return "Just now";
      }
    } catch (error) {
      console.error("Error formatting time:", error);
      return "Unknown time";
    }
  };

  if (!selectedContact) {
    return (
      <div className="chat-container">
        <div className="loading-messages">
          <h1>Loading Chat...</h1>
          <p>Please wait while we load your conversation.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="chat-container">
      {/* Chat Header */}
      <div className="chat-header">
        <button onClick={goBackToContacts} className="back-btn">
          ← Back
        </button>
        
        <div className="chat-contact-info">
          <div className="chat-avatar">
            {selectedContact.profilePicture ? (
              <img
                src={selectedContact.profilePicture}
                alt={selectedContact.name}
              />
            ) : (
              selectedContact.name?.charAt(0).toUpperCase() || '?'
            )}
          </div>
          
          <div className="chat-contact-details">
            <h2>{selectedContact.name}</h2>
            <p>{selectedContact.phoneNumber}</p>
          </div>
        </div>
        
        <div className="online-indicator" title="Auto-refreshing messages"></div>
      </div>

      {/* Error Message */}
      {error && (
        <div className="error-message">
          {error}
        </div>
      )}

      {/* Messages Area */}
      <div className="chat-messages">
        {isLoading ? (
          <div className="loading-messages">
            <p>Loading messages...</p>
          </div>
        ) : messages.length === 0 ? (
          <div className="no-messages">
            <p>No messages yet. Start the conversation!</p>
          </div>
        ) : (
          <>
            {messages.map((message) => (
              <div
                key={`${message.id}-${message.rawTimestamp}`}
                className={`message-item ${message.isMe ? 'sent' : 'received'}`}
              >
                <div className={`message-bubble ${message.isMe ? 'sent' : 'received'}`}>
                  <p className="message-text">{message.content}</p>
                  <small className="message-time">
                    {formatTime(message.timestamp)}
                  </small>
                </div>
              </div>
            ))}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      {/* Message Input */}
      <div className="chat-input">
        <form onSubmit={sendMessage} className="chat-form">
          <input
            type="text"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Type a message..."
            disabled={isSending}
            className="message-input"
          />
          <button
            type="submit"
            disabled={!newMessage.trim() || isSending}
            className="send-btn"
            title="Send message"
          >
            {isSending ? '⏳' : '➤'}
          </button>
        </form>
      </div>
    </div>
  );
}
