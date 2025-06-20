import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function Contact() {
  const [contacts, setContacts] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [showAddContact, setShowAddContact] = useState(false);
  const [phoneNumber, setPhoneNumber] = useState("");
  const [isAddingContact, setIsAddingContact] = useState(false);
  const [addContactError, setAddContactError] = useState("");
  const [addContactSuccess, setAddContactSuccess] = useState("");
  
  const navigate = useNavigate();

  useEffect(() => {
    fetchContacts();
  }, []);

  const fetchContacts = async () => {
    setIsLoading(true);
    setError("");

    try {
      // Get user ID from localStorage
      const userData = JSON.parse(localStorage.getItem('user') || '{}');
      const userId = userData.id;

      if (!userId) {
        setError("User not found. Please log in again.");
        setIsLoading(false);
        return;
      }

      console.log("Fetching contacts for user:", userId);

      const response = await fetch(`http://localhost:8080/contact/listcontacts?id=${userId}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      });

      console.log("Response status:", response.status);

      if (response.ok) {
        const data = await response.json();
        console.log("Contacts response:", data);

        if (data.data && Array.isArray(data.data)) {
          setContacts(data.data);
        } else {
          setContacts([]);
        }
      } else {
        const errorText = await response.text();
        console.log("Error response:", errorText);
        setError(errorText || "Failed to fetch contacts");
      }
    } catch (error) {
      console.error("Error fetching contacts:", error);
      setError("Something went wrong. Please try again.");
    }

    setIsLoading(false);
  };

  const addContact = async (e) => {
    e.preventDefault();
    setIsAddingContact(true);
    setAddContactError("");
    setAddContactSuccess("");

    try {
      // Get user ID from localStorage
      const userData = JSON.parse(localStorage.getItem('user') || '{}');
      const userId = userData.id;

      if (!userId) {
        setAddContactError("User not found. Please log in again.");
        setIsAddingContact(false);
        return;
      }

      if (!phoneNumber.trim()) {
        setAddContactError("Please enter a phone number.");
        setIsAddingContact(false);
        return;
      }

      console.log("Adding contact with phone:", phoneNumber);

      const response = await fetch('http://localhost:8080/contact', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_id: parseInt(userId),
          phone_number: phoneNumber.trim()
        })
      });

      console.log("Add contact response status:", response.status);

      if (response.ok) {
        const responseData = await response.text();
        console.log("Contact added successfully:", responseData);
        
        setAddContactSuccess("Contact added successfully!");
        setPhoneNumber("");
        setShowAddContact(false);
        
        // Refresh contacts list
        fetchContacts();
      } else {
        const errorText = await response.text();
        console.log("Add contact error:", errorText);
        setAddContactError(errorText || "Failed to add contact");
      }
    } catch (error) {
      console.error("Error adding contact:", error);
      setAddContactError("Something went wrong. Please try again.");
    }

    setIsAddingContact(false);
  };

  const formatPhoneNumber = (phone) => {
    if (!phone) return "No phone number";
    const cleaned = phone.replace(/\D/g, '');
    if (cleaned.length === 10) {
      return `(${cleaned.slice(0, 3)}) ${cleaned.slice(3, 6)}-${cleaned.slice(6)}`;
    }
    return phone;
  };

  const getInitials = (name) => {
    if (!name) return "?";
    return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
  };

  const handleStartChat = (contact) => {
    console.log('Starting chat with:', contact.Name, 'ID:', contact.Id);
    
    // Store the selected contact in localStorage for the chat page
    localStorage.setItem('selectedContact', JSON.stringify({
      id: contact.Id,
      name: contact.Name,
      phoneNumber: contact.PhoneNumber,
      profilePicture: contact.Profile_Picture_Url
    }));
    
    // Navigate to chat page
    navigate('/chat');
  };

  if (isLoading) {
    return (
      <div className="contacts-container">
        <div className="contacts-header">
          <h1>My Contacts</h1>
        </div>
        <div className="contacts-content">
          <p>Loading contacts...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="contacts-container">
      <div className="contacts-header">
        <h1>My Contacts</h1>
        <p>Connect with your friends and family</p>
      </div>
      
      <div className="contacts-content">
        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        {addContactSuccess && (
          <div className="success-message">
            {addContactSuccess}
          </div>
        )}

        {/* Add Contact Section */}
        <div className="add-contact-section">
          {!showAddContact ? (
            <button
              onClick={() => setShowAddContact(true)}
              className="add-contact-btn"
            >
              ➕ Add New Contact
            </button>
          ) : (
            <div className="add-contact-form">
              <h3>Add New Contact</h3>
              
              <form onSubmit={addContact}>
                <div className="input-group">
                  <label className="input-label">
                    Phone Number *
                  </label>
                  <input
                    type="tel"
                    value={phoneNumber}
                    onChange={(e) => setPhoneNumber(e.target.value)}
                    placeholder="Enter phone number (e.g., 1234567890)"
                    required
                    className="input-field"
                  />
                </div>

                {addContactError && (
                  <div className="error-message">
                    {addContactError}
                  </div>
                )}

                <div className="form-buttons">
                  <button
                    type="submit"
                    disabled={isAddingContact}
                    className="btn-primary"
                  >
                    {isAddingContact ? 'Adding...' : 'Add Contact'}
                  </button>
                  
                  <button
                    type="button"
                    onClick={() => {
                      setShowAddContact(false);
                      setPhoneNumber("");
                      setAddContactError("");
                    }}
                    className="btn-secondary"
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          )}
        </div>

        {/* Contacts List */}
        {contacts.length === 0 && !error ? (
          <div className="no-messages">
            <p>No contacts found.</p>
            <button onClick={fetchContacts} className="btn-primary">Refresh</button>
          </div>
        ) : (
          <div>
            <div className="contacts-stats">
              <span>Total contacts: {contacts.length}</span>
              <button 
                onClick={fetchContacts} 
                className="btn-secondary"
              >
                Refresh
              </button>
            </div>

            <div>
              {contacts.map((contact) => (
                <div
                  key={contact.Id}
                  className="contact-card"
                >
                  {/* Profile Picture or Initials */}
                  <div className="contact-avatar">
                    {contact.Profile_Picture_Url ? (
                      <img
                        src={contact.Profile_Picture_Url}
                        alt={contact.Name}
                        onError={(e) => {
                          e.target.style.display = 'none';
                          e.target.parentElement.innerHTML = getInitials(contact.Name);
                        }}
                      />
                    ) : (
                      getInitials(contact.Name)
                    )}
                  </div>

                  {/* Contact Info */}
                  <div className="contact-info">
                    <h3 className="contact-name">
                      {contact.Name || "Unknown"}
                    </h3>
                    
                    <div className="contact-details">
                      <div className="contact-phone">
                        📞 {formatPhoneNumber(contact.PhoneNumber)}
                      </div>
                      
                      {contact.Bio && (
                        <div className="contact-bio">
                          💬 {contact.Bio}
                        </div>
                      )}
                    </div>
                  </div>

                  {/* Chat Button */}
                  <div>
                    <button
                      className="chat-btn"
                      onClick={() => handleStartChat(contact)}
                    >
                      💬 Chat
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
