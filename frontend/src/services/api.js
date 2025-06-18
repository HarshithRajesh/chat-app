const API_BASE_URL = "http://localhost:8080/";

async function handleResponse(response) {
  if (!response.ok) {
    let errorMessage = "Something went wrong";
    try {
      const errorData = await response.json();
      if (errorData && errorData.message) {
        errorMessage = errorData.message;
      }
    } catch (jsonError) {
      errorMessage = response.statusText;
    }
    throw new Error(errorMessage);
  }
  return await response.json();
}

async function signup(userData) {
  const response = await fetch(`${API_BASE_URL}signup`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(userData),
  });
  return handleResponse(response);
}

async function login(credentials) {
  const response = await fetch(`${API_BASE_URL}Login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credentials),
  });
  return handleResponse(response);
}

export { signup, login };
