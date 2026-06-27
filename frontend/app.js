const TOKEN_KEY = "ticket_system_token";
const EMAIL_KEY = "ticket_system_email";

const alertEl = document.getElementById("alert");
const authSection = document.getElementById("auth-section");
const dashboard = document.getElementById("dashboard");
const loginForm = document.getElementById("login-form");
const registerForm = document.getElementById("register-form");
const createTicketForm = document.getElementById("create-ticket-form");
const ticketsList = document.getElementById("tickets-list");
const userEmailEl = document.getElementById("user-email");
const healthBadge = document.getElementById("health-badge");

function showAlert(message, type = "error") {
  alertEl.textContent = message;
  alertEl.className = `alert ${type}`;
  alertEl.classList.remove("hidden");
}

function hideAlert() {
  alertEl.classList.add("hidden");
}

function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

function setSession(token, email) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(EMAIL_KEY, email);
}

function clearSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(EMAIL_KEY);
}

async function api(path, options = {}) {
  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };

  const token = getToken();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(path, {
    ...options,
    headers,
  });

  let data = null;
  const contentType = response.headers.get("content-type") || "";
  if (contentType.includes("application/json")) {
    data = await response.json();
  }

  if (!response.ok) {
    const message = data?.error || `Request failed (${response.status})`;
    throw new Error(message);
  }

  return data;
}

function showAuth() {
  authSection.classList.remove("hidden");
  dashboard.classList.add("hidden");
}

function showDashboard(email) {
  authSection.classList.add("hidden");
  dashboard.classList.remove("hidden");
  userEmailEl.textContent = email;
}

function formatDate(value) {
  return new Date(value).toLocaleString();
}

function nextStatus(current) {
  if (current === "open") return "in_progress";
  if (current === "in_progress") return "closed";
  return null;
}

function renderTickets(tickets) {
  if (!tickets.length) {
    ticketsList.innerHTML = '<p class="empty">No tickets yet.</p>';
    return;
  }

  ticketsList.innerHTML = tickets
    .map((ticket) => {
      const next = nextStatus(ticket.status);
      const actionButton = next
        ? `<button type="button" data-id="${ticket.id}" data-status="${next}">Mark ${next.replace("_", " ")}</button>`
        : '<span class="empty">No further updates</span>';

      return `
        <article class="ticket">
          <div class="ticket-header">
            <h3 class="ticket-title">${escapeHtml(ticket.title)}</h3>
            <span class="status ${ticket.status}">${ticket.status.replace("_", " ")}</span>
          </div>
          <p>${escapeHtml(ticket.description || "No description")}</p>
          <p class="ticket-meta">#${ticket.id} · Created ${formatDate(ticket.created_at)}</p>
          <div class="ticket-actions">${actionButton}</div>
        </article>
      `;
    })
    .join("");

  ticketsList.querySelectorAll("button[data-id]").forEach((button) => {
    button.addEventListener("click", async () => {
      try {
        await api(`/tickets/${button.dataset.id}/status`, {
          method: "PATCH",
          body: JSON.stringify({ status: button.dataset.status }),
        });
        showAlert("Ticket status updated.", "success");
        await loadTickets();
      } catch (error) {
        showAlert(error.message);
      }
    });
  });
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

async function checkHealth() {
  try {
    const data = await api("/health");
    healthBadge.textContent = data.status === "ok" ? "API healthy" : "API issue";
  } catch {
    healthBadge.textContent = "API unreachable";
    healthBadge.style.background = "rgba(239, 68, 68, 0.15)";
    healthBadge.style.color = "#fecaca";
  }
}

async function loadTickets() {
  const tickets = await api("/tickets");
  renderTickets(tickets);
}

document.querySelectorAll(".tab").forEach((tab) => {
  tab.addEventListener("click", () => {
    document.querySelectorAll(".tab").forEach((el) => el.classList.remove("active"));
    tab.classList.add("active");

    const isLogin = tab.dataset.tab === "login";
    loginForm.classList.toggle("hidden", !isLogin);
    registerForm.classList.toggle("hidden", isLogin);
    hideAlert();
  });
});

loginForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  hideAlert();

  const formData = new FormData(loginForm);
  const email = formData.get("email");
  const password = formData.get("password");

  try {
    const data = await api("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
    setSession(data.token, email);
    showDashboard(email);
    await checkHealth();
    await loadTickets();
    showAlert("Logged in successfully.", "success");
  } catch (error) {
    showAlert(error.message);
  }
});

registerForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  hideAlert();

  const formData = new FormData(registerForm);
  const payload = {
    name: formData.get("name"),
    email: formData.get("email"),
    password: formData.get("password"),
  };

  try {
    await api("/auth/register", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    showAlert("Account created. You can log in now.", "success");
    document.querySelector('.tab[data-tab="login"]').click();
    loginForm.querySelector('[name="email"]').value = payload.email;
  } catch (error) {
    showAlert(error.message);
  }
});

createTicketForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  hideAlert();

  const formData = new FormData(createTicketForm);
  const payload = {
    title: formData.get("title"),
    description: formData.get("description"),
  };

  try {
    await api("/tickets", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    createTicketForm.reset();
    showAlert("Ticket created.", "success");
    await loadTickets();
  } catch (error) {
    showAlert(error.message);
  }
});

document.getElementById("logout-btn").addEventListener("click", () => {
  clearSession();
  showAuth();
  hideAlert();
});

document.getElementById("refresh-btn").addEventListener("click", async () => {
  try {
    await loadTickets();
    showAlert("Tickets refreshed.", "success");
  } catch (error) {
    showAlert(error.message);
  }
});

async function init() {
  const token = getToken();
  const email = localStorage.getItem(EMAIL_KEY);

  if (token && email) {
    showDashboard(email);
    try {
      await checkHealth();
      await loadTickets();
    } catch (error) {
      clearSession();
      showAuth();
      showAlert("Session expired. Please log in again.");
    }
    return;
  }

  showAuth();
  await checkHealth();
}

init();
