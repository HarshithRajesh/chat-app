import { useState, useEffect } from "react";

const Profile = () => {
  const [profileData, setProfileData] = useState({
    id: "",
    name: "",
    phone_number: "",
    bio: "",
    profile_picture_url: ""
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [imageFile, setImageFile] = useState(null);
  const [isNewUser, setIsNewUser] = useState(false);

  useEffect(() => {
    const userData = JSON.parse(localStorage.getItem('user') || '{}');
    if (userData && userData.id) {
      setProfileData(prev => ({
        ...prev,
        id: userData.id
      }));
      fetchProfile(userData.id);
    } else if (userData && userData.email) {
      // New user from signup
      setIsNewUser(true);
      setProfileData(prev => ({
        ...prev,
        id: userData.id || ""
      }));
    }
  }, []);

  const fetchProfile = async (userId) => {
    try {
      const response = await fetch(`http://localhost:8080/profile/${userId}`);
      if (response.ok) {
        const data = await response.json();
        setProfileData(data);
      } else if (response.status === 404) {
        // Profile doesn't exist yet - new user
        setIsNewUser(true);
      }
    } catch (error) {
      console.error("Error fetching profile:", error);
      setIsNewUser(true);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setProfileData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setImageFile(file);
      const previewUrl = URL.createObjectURL(file);
      setProfileData(prev => ({
        ...prev,
        profile_picture_url: previewUrl
      }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    setIsLoading(true);
    setError("");
    setSuccess("");

    try {
      // Match your Postman request format exactly
      const profilePayload = {
        id: parseInt(profileData.id),
        name: profileData.name,
        bio: profileData.bio,
        phone_number: profileData.phone_number
      };

      console.log("Sending payload:", profilePayload);

      const response = await fetch('http://localhost:8080/profile', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(profilePayload),
      });

      if (response.ok) {
        const responseData = await response.json();
        console.log("Profile updated:", responseData);
        
        setProfileData(responseData);
        setSuccess("Profile updated successfully!");
        setIsNewUser(false);
      } else {
        const errorText = await response.text();
        console.log("Error:", errorText);
        setError(errorText || "Failed to save profile");
      }
    } catch (error) {
      setError("Something went wrong. Please try again.");
      console.error("Profile save error:", error);
    }

    setIsLoading(false);
  };

  return (
    <div>
      <h1>{isNewUser ? "Complete Your Profile" : "Update Profile"}</h1>
      {isNewUser && (
        <p>Welcome! Please complete your profile to get started.</p>
      )}
      
      <form onSubmit={handleSubmit}>
        <div>
          <label>Name *</label>
          <input
            type="text"
            name="name"
            value={profileData.name}
            onChange={handleInputChange}
            required
            placeholder="Enter your full name"
          />
        </div>

        <div>
          <label>Phone Number</label>
          <input
            type="tel"
            name="phone_number"
            value={profileData.phone_number}
            onChange={handleInputChange}
            placeholder="Enter your phone number"
          />
        </div>

        <div>
          <label>Bio</label>
          <textarea
            name="bio"
            value={profileData.bio}
            onChange={handleInputChange}
            placeholder="Tell us about yourself..."
            rows="4"
          />
        </div>

        {error && <div style={{ color: 'red' }}>{error}</div>}
        {success && <div style={{ color: 'green' }}>{success}</div>}

        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Saving...' : (isNewUser ? 'Create Profile' : 'Update Profile')}
        </button>
      </form>
    </div>
  );
};

export default Profile;