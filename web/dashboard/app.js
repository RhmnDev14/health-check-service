const API_BASE = "/api";

// Auto-refresh interval (10 seconds)
const REFRESH_INTERVAL = 10000;

// DOM Ready
document.addEventListener("DOMContentLoaded", () => {
  loadDashboardStats();
  loadRecipients();
  loadRecentLogs();

  // Auto refresh
  setInterval(() => {
    loadDashboardStats();
    loadRecentLogs();
  }, REFRESH_INTERVAL);
});

// Load dashboard stats and services
async function loadDashboardStats() {
  try {
    const response = await fetch(`${API_BASE}/dashboard/stats`);
    const data = await response.json();

    document.getElementById("totalServices").textContent =
      data.total_services || 0;
    document.getElementById("upCount").textContent = data.up_count || 0;
    document.getElementById("downCount").textContent = data.down_count || 0;

    renderServices(data.services || []);
  } catch (error) {
    console.error("Failed to load dashboard stats:", error);
  }
}

// Render services
function renderServices(services) {
  const grid = document.getElementById("servicesGrid");

  if (services.length === 0) {
    grid.innerHTML =
      '<div class="empty-state">No services configured. Add your first service!</div>';
    return;
  }

  grid.innerHTML = services
    .map((service) => {
      const statusClass = service.status.toLowerCase();
      return `
            <div class="service-card ${statusClass}">
                <div class="service-name">${escapeHtml(service.service_name)}</div>
                <div class="service-url">${escapeHtml(service.url)}</div>
                <div class="service-status">
                    <span class="status-badge ${statusClass}">${service.status}</span>
                    ${service.response_time > 0 ? `<span class="service-meta">${service.response_time}ms</span>` : ""}
                </div>
                ${service.error_message ? `<div class="service-meta" style="color: var(--status-down);">${escapeHtml(service.error_message)}</div>` : ""}
                <div class="service-actions">
                    <button class="btn btn-sm btn-secondary" onclick="toggleService('${service.service_id}', ${!service.is_active})">
                        ${service.is_active ? "⏸ Pause" : "▶ Resume"}
                    </button>
                    <button class="btn btn-sm btn-danger" onclick="deleteService('${service.service_id}')">🗑 Delete</button>
                </div>
            </div>
        `;
    })
    .join("");
}

// Load recipients
async function loadRecipients() {
  try {
    const response = await fetch(`${API_BASE}/recipients`);
    const recipients = await response.json();
    renderRecipients(recipients || []);
  } catch (error) {
    console.error("Failed to load recipients:", error);
  }
}

// Render recipients
function renderRecipients(recipients) {
  const list = document.getElementById("recipientsList");

  if (recipients.length === 0) {
    list.innerHTML =
      '<div class="empty-state">No recipients configured. Add recipients to receive notifications!</div>';
    return;
  }

  list.innerHTML = recipients
    .map(
      (recipient) => `
        <div class="recipient-chip ${recipient.is_active ? "" : "inactive"}">
            <span class="recipient-name">${escapeHtml(recipient.name)}</span>
            <span class="recipient-phone">${escapeHtml(recipient.phone)}</span>
            <button class="btn btn-sm btn-secondary" onclick="toggleRecipient('${recipient.id}', ${!recipient.is_active})">
                ${recipient.is_active ? "⏸" : "▶"}
            </button>
            <button class="btn btn-sm btn-danger" onclick="deleteRecipient('${recipient.id}')">🗑</button>
        </div>
    `,
    )
    .join("");
}

// Load recent logs
async function loadRecentLogs() {
  try {
    const response = await fetch(`${API_BASE}/health-logs?limit=20`);
    const logs = await response.json();
    renderLogs(logs || []);
  } catch (error) {
    console.error("Failed to load logs:", error);
  }
}

// Render logs
function renderLogs(logs) {
  const tbody = document.getElementById("logsTableBody");

  if (logs.length === 0) {
    tbody.innerHTML =
      '<tr><td colspan="5" class="empty-state">No health check logs yet.</td></tr>';
    return;
  }

  tbody.innerHTML = logs
    .map((log) => {
      const checkedAt = new Date(log.checked_at).toLocaleString();
      const statusClass = log.status.toLowerCase();
      return `
            <tr>
                <td>${checkedAt}</td>
                <td>${escapeHtml(log.service_name)}</td>
                <td><span class="status-badge ${statusClass}">${log.status}</span></td>
                <td>${log.response_time_ms}ms</td>
                <td>${log.error_message ? escapeHtml(log.error_message) : "-"}</td>
            </tr>
        `;
    })
    .join("");
}

// Add service
async function addService(event) {
  event.preventDefault();
  const form = event.target;
  const formData = new FormData(form);
  const data = {
    id: formData.get("id"),
    name: formData.get("name"),
    url: formData.get("url"),
    check_interval: parseInt(formData.get("check_interval")) || 60,
    retry_interval: parseInt(formData.get("retry_interval")) || 180,
  };

  try {
    const response = await fetch(`${API_BASE}/services`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });

    if (response.ok) {
      form.reset();
      closeModal("addServiceModal");
      loadDashboardStats();
    } else {
      const error = await response.json();
      alert("Failed to add service: " + (error.error || "Unknown error"));
    }
  } catch (error) {
    alert("Failed to add service: " + error.message);
  }
}

// Toggle service active status
async function toggleService(id, isActive) {
  try {
    await fetch(`${API_BASE}/services/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ is_active: isActive }),
    });
    loadDashboardStats();
  } catch (error) {
    alert("Failed to update service: " + error.message);
  }
}

// Delete service
async function deleteService(id) {
  if (!confirm("Are you sure you want to delete this service?")) return;

  try {
    await fetch(`${API_BASE}/services/${id}`, { method: "DELETE" });
    loadDashboardStats();
  } catch (error) {
    alert("Failed to delete service: " + error.message);
  }
}

// Add recipient
async function addRecipient(event) {
  event.preventDefault();
  const form = event.target;
  const formData = new FormData(form);
  const data = {
    name: formData.get("name"),
    phone: formData.get("phone"),
  };

  try {
    const response = await fetch(`${API_BASE}/recipients`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });

    if (response.ok) {
      form.reset();
      closeModal("addRecipientModal");
      loadRecipients();
    } else {
      const error = await response.json();
      alert("Failed to add recipient: " + (error.error || "Unknown error"));
    }
  } catch (error) {
    alert("Failed to add recipient: " + error.message);
  }
}

// Toggle recipient active status
async function toggleRecipient(id, isActive) {
  try {
    await fetch(`${API_BASE}/recipients/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ is_active: isActive }),
    });
    loadRecipients();
  } catch (error) {
    alert("Failed to update recipient: " + error.message);
  }
}

// Delete recipient
async function deleteRecipient(id) {
  if (!confirm("Are you sure you want to delete this recipient?")) return;

  try {
    await fetch(`${API_BASE}/recipients/${id}`, { method: "DELETE" });
    loadRecipients();
  } catch (error) {
    alert("Failed to delete recipient: " + error.message);
  }
}

// Modal functions
function showAddServiceModal() {
  document.getElementById("addServiceModal").classList.add("active");
}

function showAddRecipientModal() {
  document.getElementById("addRecipientModal").classList.add("active");
}

function closeModal(id) {
  document.getElementById(id).classList.remove("active");
}

// Close modal on outside click
document.querySelectorAll(".modal").forEach((modal) => {
  modal.addEventListener("click", (e) => {
    if (e.target === modal) {
      modal.classList.remove("active");
    }
  });
});

// Escape HTML to prevent XSS
function escapeHtml(text) {
  if (!text) return "";
  const div = document.createElement("div");
  div.textContent = text;
  return div.innerHTML;
}
