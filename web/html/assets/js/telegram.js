const baseUrl = document.body.dataset.baseUrl || "/";
const apiBase = `${baseUrl}api/`;

const state = {
    config: null,
    tariffs: [],
    broadcasts: [],
    promos: [],
    conversations: [],
    selectedTariffId: null,
    selectedBroadcastId: null,
    activeConversation: null,
};

const el = {
    toast: document.getElementById("toast"),
    loading: document.getElementById("loading-backdrop"),
    configForm: document.getElementById("config-form"),
    configEnabled: document.getElementById("config-enabled"),
    botTokenInput: document.getElementById("config-bot-token"),
    botTokenMask: document.getElementById("bot-token-mask"),
    webhookDomain: document.getElementById("config-webhook-domain"),
    webhookSecret: document.getElementById("config-webhook-secret"),
    shopId: document.getElementById("config-shop-id"),
    secretKey: document.getElementById("config-secret-key"),
    successUrl: document.getElementById("config-success-url"),
    failureUrl: document.getElementById("config-failure-url"),
    miniAppUrl: document.getElementById("config-mini-app"),
    downloadLinksBody: document.getElementById("download-links-body"),
    addDownloadLink: document.getElementById("add-download-link"),
    tariffTableBody: document.getElementById("tariff-table-body"),
    createTariff: document.getElementById("create-tariff"),
    refreshTariffs: document.getElementById("refresh-tariffs"),
    tariffForm: document.getElementById("tariff-form"),
    tariffId: document.getElementById("tariff-id"),
    tariffTitle: document.getElementById("tariff-title"),
    tariffDescription: document.getElementById("tariff-description"),
    tariffPrice: document.getElementById("tariff-price"),
    tariffCurrency: document.getElementById("tariff-currency"),
    tariffDuration: document.getElementById("tariff-duration"),
    tariffSort: document.getElementById("tariff-sort"),
    tariffActive: document.getElementById("tariff-active"),
    buttonsSubtitle: document.getElementById("buttons-subtitle"),
    buttonTableBody: document.getElementById("button-table-body"),
    buttonForm: document.getElementById("button-form"),
    buttonId: document.getElementById("button-id"),
    buttonTariffId: document.getElementById("button-tariff-id"),
    buttonLabel: document.getElementById("button-label"),
    buttonAction: document.getElementById("button-action"),
    buttonPayload: document.getElementById("button-payload"),
    buttonSort: document.getElementById("button-sort"),
    cancelButtonEdit: document.getElementById("cancel-button-edit"),
    addButton: document.getElementById("add-button"),
    broadcastTableBody: document.getElementById("broadcast-table-body"),
    broadcastForm: document.getElementById("broadcast-form"),
    broadcastId: document.getElementById("broadcast-id"),
    broadcastTitle: document.getElementById("broadcast-title"),
    broadcastBody: document.getElementById("broadcast-body"),
    broadcastEditable: document.getElementById("broadcast-editable"),
    broadcastAllUsers: document.getElementById("broadcast-all-users"),
    broadcastIncludeNever: document.getElementById("broadcast-include-never"),
    broadcastIncludeExpired: document.getElementById("broadcast-include-expired"),
    broadcastTariffOptions: document.getElementById("broadcast-tariff-options"),
    broadcastCreate: document.getElementById("create-broadcast"),
    broadcastSend: document.getElementById("broadcast-send"),
    broadcastRefresh: document.getElementById("refresh-broadcasts"),
    broadcastEditSent: document.getElementById("broadcast-edit-sent"),
    broadcastStatus: document.getElementById("broadcast-status"),
    promoTableBody: document.getElementById("promo-table-body"),
    promoForm: document.getElementById("promo-form"),
    promoId: document.getElementById("promo-id"),
    promoCode: document.getElementById("promo-code"),
    promoDescription: document.getElementById("promo-description"),
    promoDiscount: document.getElementById("promo-discount"),
    promoFreeDays: document.getElementById("promo-free-days"),
    promoMaxUses: document.getElementById("promo-max-uses"),
    promoActive: document.getElementById("promo-active"),
    promoExpiry: document.getElementById("promo-expiry"),
    promoNoExpiry: document.getElementById("promo-no-expiry"),
    conversationList: document.getElementById("conversation-list"),
    conversationHeader: document.getElementById("conversation-header"),
    conversationMessages: document.getElementById("conversation-messages"),
    conversationEmpty: document.getElementById("conversation-empty"),
    conversationReplyForm: document.getElementById("conversation-reply-form"),
    conversationReply: document.getElementById("conversation-reply"),
    refreshConversations: document.getElementById("refresh-conversations"),
};

function showToast(message, isError = false) {
    if (!message) {
        el.toast.classList.remove("visible", "error");
        el.toast.textContent = "";
        return;
    }
    el.toast.textContent = message;
    el.toast.classList.toggle("error", isError);
    requestAnimationFrame(() => {
        el.toast.classList.add("visible");
    });
    clearTimeout(showToast.timer);
    showToast.timer = setTimeout(() => {
        el.toast.classList.remove("visible", "error");
    }, 4000);
}

function setLoading(isLoading) {
    el.loading.toggleAttribute("hidden", !isLoading);
}

async function request(path, options = {}) {
    const { headers, ...rest } = options;
    const response = await fetch(`${apiBase}${path}`, {
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            ...(headers || {}),
        },
        ...rest,
    });

    if (response.status === 401) {
        window.location.href = `${baseUrl}login`;
        return Promise.reject(new Error("Требуется авторизация"));
    }

    const data = await response.json().catch(() => ({}));
    if (!data.success) {
        throw new Error(data.msg || "Не удалось выполнить запрос");
    }
    return data.obj;
}

async function loadState(preserveSelection = true) {
    setLoading(true);
    try {
        const obj = await request("telegramState", { method: "GET" });
        state.config = obj?.config || null;
        state.tariffs = Array.isArray(obj?.tariffs) ? obj.tariffs : [];
        state.broadcasts = Array.isArray(obj?.broadcasts) ? obj.broadcasts : [];
        state.promos = Array.isArray(obj?.promoCodes) ? obj.promoCodes : [];
        state.conversations = Array.isArray(obj?.conversations) ? obj.conversations : [];
        if (preserveSelection && state.selectedTariffId) {
            const exists = state.tariffs.some((t) => t.id === state.selectedTariffId);
            if (!exists) {
                state.selectedTariffId = state.tariffs[0]?.id || null;
            }
        } else if (!state.selectedTariffId) {
            state.selectedTariffId = state.tariffs[0]?.id || null;
        }
        if (!state.selectedBroadcastId || !state.broadcasts.some((b) => b.id === state.selectedBroadcastId)) {
            state.selectedBroadcastId = state.broadcasts[0]?.id || null;
        }
        renderConfig();
        renderTariffs();
        renderButtons();
        renderBroadcasts();
        renderBroadcastAudienceOptions();
        renderPromos();
        renderConversations();
        if (state.activeConversation) {
            await loadConversation(state.activeConversation.user.id, false);
        }
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
}

function renderConfig() {
    const cfg = state.config || {};
    el.configEnabled.checked = Boolean(cfg.enabled);
    el.webhookDomain.value = cfg.webhookDomain || "";
    el.webhookSecret.value = cfg.webhookSecret || "";
    el.shopId.value = cfg.yooKassaShopId || "";
    el.secretKey.value = "";
    el.successUrl.value = cfg.successRedirectUrl || "";
    el.failureUrl.value = cfg.failureRedirectUrl || "";
    el.miniAppUrl.value = cfg.miniAppUrl || "";
    el.botTokenInput.value = "";
    if (cfg.botTokenMasked) {
        el.botTokenMask.textContent = `Сохранён токен: ${cfg.botTokenMasked}`;
    } else {
        el.botTokenMask.textContent = "Токен не задан.";
    }

    renderDownloadLinks(cfg.downloadLinks || {});
}

function renderDownloadLinks(map) {
    el.downloadLinksBody.innerHTML = "";
    const entries = Object.entries(map || {});
    if (entries.length === 0) {
        addDownloadLinkRow();
        return;
    }
    for (const [platform, url] of entries) {
        addDownloadLinkRow(platform, url);
    }
}

function addDownloadLinkRow(platform = "", url = "") {
    const row = document.createElement("tr");
    row.innerHTML = `
        <td><input type="text" class="input" value="${escapeHtml(platform)}" placeholder="Windows" /></td>
        <td><input type="url" class="input" value="${escapeHtml(url)}" placeholder="https://..." /></td>
        <td class="actions"><button type="button" class="btn subtle remove-link">Удалить</button></td>
    `;
    const removeBtn = row.querySelector(".remove-link");
    removeBtn.addEventListener("click", (event) => {
        event.stopPropagation();
        row.remove();
        if (!el.downloadLinksBody.children.length) {
            addDownloadLinkRow();
        }
    });
    el.downloadLinksBody.appendChild(row);
}

function collectDownloadLinks() {
    const links = {};
    for (const row of el.downloadLinksBody.querySelectorAll("tr")) {
        const inputs = row.querySelectorAll("input");
        const platform = inputs[0]?.value.trim();
        const url = inputs[1]?.value.trim();
        if (platform && url) {
            links[platform] = url;
        }
    }
    return links;
}

function formatPrice(priceMinor, currency) {
    if (typeof priceMinor !== "number") {
        return "—";
    }
    const major = priceMinor / 100;
    return `${major.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })} ${currency || ""}`.trim();
}

function formatDateTime(value) {
    if (!value) {
        return "";
    }
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
        return "";
    }
    return date.toLocaleString();
}

function describeAudience(audience) {
    if (!audience || audience.allUsers) {
        return "Все пользователи";
    }
    const parts = [];
    if (Array.isArray(audience.tariffIds) && audience.tariffIds.length) {
        const names = audience.tariffIds
            .map((id) => state.tariffs.find((t) => t.id === id)?.title || `Тариф #${id}`)
            .filter(Boolean);
        if (names.length) {
            parts.push(`Покупатели: ${names.join(", ")}`);
        }
    }
    if (audience.includeNeverSubscribed) {
        parts.push("Никогда не оформлявшие");
    }
    if (audience.includeExpired) {
        parts.push("С истёкшей подпиской");
    }
    return parts.join(", ") || "Выбранные пользователи";
}

function getBroadcastTariffSelection() {
    if (!el.broadcastTariffOptions) {
        return [];
    }
    return Array.from(
        el.broadcastTariffOptions.querySelectorAll('input[type="checkbox"]:checked')
    ).map((input) => Number(input.value)).filter((value) => !Number.isNaN(value));
}

function renderBroadcastAudienceOptions(selectedTariffs) {
    if (!el.broadcastTariffOptions) {
        return;
    }
    const selected = new Set(
        (Array.isArray(selectedTariffs) ? selectedTariffs : getBroadcastTariffSelection()).map((id) => Number(id))
    );
    el.broadcastTariffOptions.innerHTML = "";
    if (!state.tariffs.length) {
        const hint = document.createElement("p");
        hint.className = "hint";
        hint.textContent = "Создайте тарифы, чтобы таргетировать рассылки.";
        el.broadcastTariffOptions.appendChild(hint);
        return;
    }
    for (const tariff of state.tariffs) {
        const id = Number(tariff.id);
        const label = document.createElement("label");
        label.className = "checkbox";
        const input = document.createElement("input");
        input.type = "checkbox";
        input.value = String(id);
        input.checked = selected.has(id);
        label.appendChild(input);
        const span = document.createElement("span");
        span.textContent = tariff.title || `Тариф #${id}`;
        label.appendChild(span);
        el.broadcastTariffOptions.appendChild(label);
    }
    updateBroadcastAudienceControls();
}

function updateBroadcastAudienceControls() {
    if (!el.broadcastTariffOptions) {
        return;
    }
    const disable = el.broadcastAllUsers?.checked;
    for (const input of el.broadcastTariffOptions.querySelectorAll('input[type="checkbox"]')) {
        input.disabled = disable;
    }
}

function renderTariffs() {
    el.tariffTableBody.innerHTML = "";
    if (!state.tariffs.length) {
        const emptyRow = document.createElement("tr");
        const cell = document.createElement("td");
        cell.colSpan = 5;
        cell.textContent = "Тарифы не созданы.";
        cell.classList.add("hint");
        emptyRow.appendChild(cell);
        el.tariffTableBody.appendChild(emptyRow);
        resetTariffForm();
        return;
    }

    for (const tariff of state.tariffs) {
        const row = document.createElement("tr");
        row.dataset.id = String(tariff.id);
        row.innerHTML = `
            <td>${escapeHtml(tariff.title || "Без названия")}</td>
            <td>${formatPrice(tariff.priceMinor, tariff.currency)}</td>
            <td>${tariff.durationDays ? `${tariff.durationDays} д.` : "Бессрочно"}</td>
            <td>
                <span class="badge ${tariff.active ? "success" : "muted"}">
                    ${tariff.active ? "Активен" : "Выключен"}
                </span>
            </td>
            <td class="actions"></td>
        `;
        if (state.selectedTariffId === tariff.id) {
            row.classList.add("selected");
        }
        const actionsCell = row.querySelector(".actions");
        const editBtn = document.createElement("button");
        editBtn.type = "button";
        editBtn.className = "btn subtle";
        editBtn.textContent = "Изменить";
        editBtn.dataset.action = "edit";
        editBtn.dataset.id = String(tariff.id);
        const deleteBtn = document.createElement("button");
        deleteBtn.type = "button";
        deleteBtn.className = "btn subtle";
        deleteBtn.textContent = "Удалить";
        deleteBtn.dataset.action = "delete";
        deleteBtn.dataset.id = String(tariff.id);
        actionsCell.append(editBtn, deleteBtn);
        el.tariffTableBody.appendChild(row);
    }

    const selected = state.tariffs.find((t) => t.id === state.selectedTariffId);
    if (selected) {
        populateTariffForm(selected);
    } else {
        resetTariffForm();
    }
}

function resetTariffForm() {
    el.tariffForm.reset();
    el.tariffId.value = "";
    el.tariffPrice.value = "";
    el.tariffSort.value = "0";
    el.tariffActive.checked = true;
}

function populateTariffForm(tariff) {
    el.tariffId.value = tariff.id || "";
    el.tariffTitle.value = tariff.title || "";
    el.tariffDescription.value = tariff.description || "";
    el.tariffPrice.value = tariff.priceMinor ? (tariff.priceMinor / 100).toFixed(2) : "";
    el.tariffCurrency.value = tariff.currency || "";
    el.tariffDuration.value = typeof tariff.durationDays === "number" ? tariff.durationDays : "";
    el.tariffSort.value = typeof tariff.sortOrder === "number" ? tariff.sortOrder : "0";
    el.tariffActive.checked = Boolean(tariff.active);
}

function selectTariff(id) {
    state.selectedTariffId = id ? Number(id) : null;
    for (const row of el.tariffTableBody.querySelectorAll("tr")) {
        row.classList.toggle("selected", Number(row.dataset.id) === state.selectedTariffId);
    }
    const selected = state.tariffs.find((t) => t.id === state.selectedTariffId);
    if (selected) {
        populateTariffForm(selected);
    }
    renderButtons();
}

function renderButtons() {
    el.buttonTableBody.innerHTML = "";
    const tariff = state.tariffs.find((t) => t.id === state.selectedTariffId);
    if (!tariff) {
        el.buttonsSubtitle.textContent = "Выберите тариф, чтобы управлять кнопками.";
        el.buttonForm.hidden = true;
        return;
    }
    el.buttonsSubtitle.textContent = `Кнопки тарифа «${tariff.title || "Без названия"}»`;
    el.buttonForm.hidden = true;
    resetButtonForm();

    if (!Array.isArray(tariff.buttons) || !tariff.buttons.length) {
        const row = document.createElement("tr");
        const cell = document.createElement("td");
        cell.colSpan = 4;
        cell.textContent = "Кнопки не настроены.";
        cell.classList.add("hint");
        row.appendChild(cell);
        el.buttonTableBody.appendChild(row);
    } else {
        for (const button of tariff.buttons) {
            const row = document.createElement("tr");
            row.dataset.id = String(button.id);
            row.innerHTML = `
                <td>${escapeHtml(button.label)}</td>
                <td>${escapeHtml(button.action)}</td>
                <td>${escapeHtml(button.payload || "")}</td>
                <td class="actions"></td>
            `;
            const actionsCell = row.querySelector(".actions");
            const editBtn = document.createElement("button");
            editBtn.type = "button";
            editBtn.className = "btn subtle";
            editBtn.textContent = "Изменить";
            editBtn.dataset.action = "edit";
            editBtn.dataset.id = String(button.id);
            const deleteBtn = document.createElement("button");
            deleteBtn.type = "button";
            deleteBtn.className = "btn subtle";
            deleteBtn.textContent = "Удалить";
            deleteBtn.dataset.action = "delete";
            deleteBtn.dataset.id = String(button.id);
            actionsCell.append(editBtn, deleteBtn);
            el.buttonTableBody.appendChild(row);
        }
    }
}

function escapeHtml(value = "") {
    return String(value ?? "").replace(/[&<>"']/g, (ch) => ({
        "&": "&amp;",
        "<": "&lt;",
        ">": "&gt;",
        '"': "&quot;",
        "'": "&#39;",
    })[ch]);
}

function resetButtonForm() {
    el.buttonForm.reset();
    el.buttonId.value = "";
    el.buttonTariffId.value = state.selectedTariffId || "";
    el.buttonSort.value = "0";
}

function openButtonForm(button) {
    if (!state.selectedTariffId) {
        showToast("Выберите тариф", true);
        return;
    }
    el.buttonForm.hidden = false;
    if (button) {
        el.buttonId.value = button.id || "";
        el.buttonTariffId.value = button.tariffId || state.selectedTariffId;
        el.buttonLabel.value = button.label || "";
        el.buttonAction.value = button.action || "";
        el.buttonPayload.value = button.payload || "";
        el.buttonSort.value = typeof button.sortOrder === "number" ? button.sortOrder : "0";
    } else {
        resetButtonForm();
    }
    el.buttonLabel.focus();
}

function closeButtonForm() {
    el.buttonForm.hidden = true;
    resetButtonForm();
}

function renderBroadcasts() {
    if (!el.broadcastTableBody) {
        return;
    }
    el.broadcastTableBody.innerHTML = "";
    if (!Array.isArray(state.broadcasts) || !state.broadcasts.length) {
        const row = document.createElement("tr");
        const cell = document.createElement("td");
        cell.colSpan = 5;
        cell.textContent = "Рассылки ещё не созданы.";
        cell.classList.add("hint");
        row.appendChild(cell);
        el.broadcastTableBody.appendChild(row);
        resetBroadcastForm();
        return;
    }
    for (const broadcast of state.broadcasts) {
        const row = document.createElement("tr");
        row.dataset.id = String(broadcast.id);
        row.classList.toggle("selected", broadcast.id === state.selectedBroadcastId);
        const deliveries = broadcast.deliveries || 0;
        const success = broadcast.success || 0;
        const failed = broadcast.failed || 0;
        row.innerHTML = `
            <td>${escapeHtml(broadcast.title || "Без названия")}</td>
            <td>${escapeHtml(describeAudience(broadcast.audience))}</td>
            <td>
                <span class="badge ${broadcast.status === "sent" ? "success" : "muted"}">
                    ${broadcast.status === "sent" ? "Отправлено" : "Черновик"}
                </span>
            </td>
            <td>${deliveries ? `${success}/${deliveries}` : "—"}</td>
            <td class="actions"></td>
        `;
        const actionsCell = row.querySelector(".actions");
        const editBtn = document.createElement("button");
        editBtn.type = "button";
        editBtn.className = "btn subtle";
        editBtn.textContent = "Редактировать";
        editBtn.dataset.action = "edit";
        editBtn.dataset.id = String(broadcast.id);
        actionsCell.appendChild(editBtn);
        if (broadcast.status !== "sent") {
            const sendBtn = document.createElement("button");
            sendBtn.type = "button";
            sendBtn.className = "btn subtle";
            sendBtn.textContent = "Отправить";
            sendBtn.dataset.action = "send";
            sendBtn.dataset.id = String(broadcast.id);
            actionsCell.appendChild(sendBtn);
        } else if (broadcast.editable) {
            const editMessageBtn = document.createElement("button");
            editMessageBtn.type = "button";
            editMessageBtn.className = "btn subtle";
            editMessageBtn.textContent = "Обновить текст";
            editMessageBtn.dataset.action = "edit-sent";
            editMessageBtn.dataset.id = String(broadcast.id);
            actionsCell.appendChild(editMessageBtn);
        }
        const deleteBtn = document.createElement("button");
        deleteBtn.type = "button";
        deleteBtn.className = "btn subtle";
        deleteBtn.textContent = "Удалить";
        deleteBtn.dataset.action = "delete";
        deleteBtn.dataset.id = String(broadcast.id);
        deleteBtn.disabled = broadcast.status !== "draft";
        actionsCell.appendChild(deleteBtn);
        el.broadcastTableBody.appendChild(row);
    }
    const selected = state.broadcasts.find((b) => b.id === state.selectedBroadcastId);
    if (selected) {
        populateBroadcastForm(selected);
    } else {
        resetBroadcastForm();
    }
}

function resetBroadcastForm() {
    if (!el.broadcastForm) {
        return;
    }
    el.broadcastForm.reset();
    if (el.broadcastId) {
        el.broadcastId.value = "";
    }
    if (el.broadcastTitle) {
        el.broadcastTitle.value = "";
    }
    if (el.broadcastBody) {
        el.broadcastBody.value = "";
    }
    if (el.broadcastEditable) {
        el.broadcastEditable.checked = true;
    }
    if (el.broadcastAllUsers) {
        el.broadcastAllUsers.checked = true;
    }
    if (el.broadcastIncludeNever) {
        el.broadcastIncludeNever.checked = false;
    }
    if (el.broadcastIncludeExpired) {
        el.broadcastIncludeExpired.checked = false;
    }
    renderBroadcastAudienceOptions([]);
    updateBroadcastActionButtons(null);
    if (el.broadcastStatus) {
        el.broadcastStatus.textContent = "Создайте новую рассылку или выберите существующую.";
    }
}

function updateBroadcastActionButtons(broadcast) {
    if (!el.broadcastSend || !el.broadcastEditSent) {
        return;
    }
    if (!broadcast || !broadcast.id) {
        el.broadcastSend.disabled = true;
        el.broadcastSend.hidden = true;
        el.broadcastEditSent.hidden = true;
        if (el.broadcastStatus) {
            el.broadcastStatus.textContent = "Черновик не сохранён.";
        }
        return;
    }
    const isDraft = broadcast.status !== "sent";
    el.broadcastSend.hidden = false;
    el.broadcastSend.disabled = !isDraft;
    el.broadcastSend.dataset.id = String(broadcast.id);
    el.broadcastEditSent.hidden = !(broadcast.status === "sent" && broadcast.editable);
    el.broadcastEditSent.disabled = !(broadcast.status === "sent" && broadcast.editable);
    el.broadcastEditSent.dataset.id = String(broadcast.id);
    if (el.broadcastStatus) {
        const pieces = [];
        if (broadcast.status === "sent") {
            pieces.push(`Отправлено ${broadcast.sentAt ? formatDateTime(broadcast.sentAt) : ""}`.trim());
            pieces.push(`Доставлено: ${broadcast.success || 0}`);
            if (broadcast.failed) {
                pieces.push(`Ошибки: ${broadcast.failed}`);
            }
            if (broadcast.editable) {
                pieces.push("Можно обновить текст");
            }
        } else {
            pieces.push("Черновик не отправлен");
        }
        el.broadcastStatus.textContent = pieces.filter(Boolean).join(" • ");
    }
}

function populateBroadcastForm(broadcast) {
    if (!el.broadcastForm || !broadcast) {
        resetBroadcastForm();
        return;
    }
    state.selectedBroadcastId = broadcast.id || null;
    if (el.broadcastId) {
        el.broadcastId.value = broadcast.id || "";
    }
    if (el.broadcastTitle) {
        el.broadcastTitle.value = broadcast.title || "";
    }
    if (el.broadcastBody) {
        el.broadcastBody.value = broadcast.body || "";
    }
    if (el.broadcastEditable) {
        el.broadcastEditable.checked = Boolean(broadcast.editable);
    }
    const audience = broadcast.audience || {};
    if (el.broadcastAllUsers) {
        el.broadcastAllUsers.checked = Boolean(audience.allUsers);
    }
    if (el.broadcastIncludeNever) {
        el.broadcastIncludeNever.checked = Boolean(audience.includeNeverSubscribed);
    }
    if (el.broadcastIncludeExpired) {
        el.broadcastIncludeExpired.checked = Boolean(audience.includeExpired);
    }
    renderBroadcastAudienceOptions(audience.tariffIds || []);
    updateBroadcastActionButtons(broadcast);
}

function collectBroadcastPayload() {
    const tariffs = getBroadcastTariffSelection();
    return {
        id: el.broadcastId?.value ? Number(el.broadcastId.value) : 0,
        title: el.broadcastTitle?.value.trim() || "",
        body: el.broadcastBody?.value.trim() || "",
        editable: Boolean(el.broadcastEditable?.checked),
        audience: {
            allUsers: Boolean(el.broadcastAllUsers?.checked),
            tariffIds: tariffs,
            includeNeverSubscribed: Boolean(el.broadcastIncludeNever?.checked),
            includeExpired: Boolean(el.broadcastIncludeExpired?.checked),
        },
    };
}

function renderPromos() {
    if (!el.promoTableBody) {
        return;
    }
    el.promoTableBody.innerHTML = "";
    if (!Array.isArray(state.promos) || !state.promos.length) {
        const row = document.createElement("tr");
        const cell = document.createElement("td");
        cell.colSpan = 6;
        cell.className = "hint";
        cell.textContent = "Промокоды ещё не созданы.";
        row.appendChild(cell);
        el.promoTableBody.appendChild(row);
        resetPromoForm();
        return;
    }
    for (const promo of state.promos) {
        const row = document.createElement("tr");
        row.dataset.id = String(promo.id);
        row.innerHTML = `
            <td>${escapeHtml(promo.code)}</td>
            <td>${promo.discountPercent ? `${promo.discountPercent}%` : "—"}</td>
            <td>${promo.freeDays ? `${promo.freeDays} д.` : "—"}</td>
            <td>${promo.maxUses ? `${promo.usedCount}/${promo.maxUses}` : promo.usedCount || 0}</td>
            <td>${promo.expiresAt ? formatDateTime(promo.expiresAt) : "Без срока"}</td>
            <td class="actions"></td>
        `;
        const actionsCell = row.querySelector(".actions");
        const editBtn = document.createElement("button");
        editBtn.type = "button";
        editBtn.className = "btn subtle";
        editBtn.textContent = "Изменить";
        editBtn.dataset.action = "edit";
        editBtn.dataset.id = String(promo.id);
        const deleteBtn = document.createElement("button");
        deleteBtn.type = "button";
        deleteBtn.className = "btn subtle";
        deleteBtn.textContent = "Удалить";
        deleteBtn.dataset.action = "delete";
        deleteBtn.dataset.id = String(promo.id);
        actionsCell.append(editBtn, deleteBtn);
        el.promoTableBody.appendChild(row);
    }
}

function resetPromoForm() {
    if (!el.promoForm) {
        return;
    }
    el.promoForm.reset();
    el.promoId.value = "";
    el.promoDiscount.value = "0";
    el.promoFreeDays.value = "0";
    el.promoMaxUses.value = "0";
    el.promoActive.checked = true;
    el.promoNoExpiry.checked = true;
    el.promoExpiry.value = "";
    el.promoExpiry.disabled = true;
}

function populatePromoForm(promo) {
    if (!promo) {
        resetPromoForm();
        return;
    }
    el.promoId.value = promo.id || "";
    el.promoCode.value = promo.code || "";
    el.promoDescription.value = promo.description || "";
    el.promoDiscount.value = typeof promo.discountPercent === "number" ? promo.discountPercent : 0;
    el.promoFreeDays.value = typeof promo.freeDays === "number" ? promo.freeDays : 0;
    el.promoMaxUses.value = typeof promo.maxUses === "number" ? promo.maxUses : 0;
    el.promoActive.checked = Boolean(promo.active);
    if (promo.expiresAt) {
        el.promoNoExpiry.checked = false;
        const date = new Date(promo.expiresAt);
        if (!Number.isNaN(date.getTime())) {
            el.promoExpiry.value = date.toISOString().slice(0, 10);
        }
        el.promoExpiry.disabled = false;
    } else {
        el.promoNoExpiry.checked = true;
        el.promoExpiry.value = "";
        el.promoExpiry.disabled = true;
    }
}

function collectPromoPayload() {
    const noExpiry = el.promoNoExpiry.checked;
    let expiresAt = null;
    if (!noExpiry && el.promoExpiry.value) {
        const date = new Date(el.promoExpiry.value);
        if (!Number.isNaN(date.getTime())) {
            expiresAt = date.toISOString();
        }
    }
    return {
        id: el.promoId.value ? Number(el.promoId.value) : 0,
        code: el.promoCode.value.trim(),
        description: el.promoDescription.value.trim(),
        discountPercent: Number(el.promoDiscount.value || 0),
        freeDays: Number(el.promoFreeDays.value || 0),
        maxUses: Number(el.promoMaxUses.value || 0),
        active: el.promoActive.checked,
        noExpiry,
        expiresAt,
    };
}

function renderConversations() {
    if (!el.conversationList) {
        return;
    }
    el.conversationList.innerHTML = "";
    if (!Array.isArray(state.conversations) || !state.conversations.length) {
        const empty = document.createElement("p");
        empty.className = "hint";
        empty.textContent = "Сообщений от пользователей ещё не поступало.";
        el.conversationList.appendChild(empty);
        renderConversationThread(null);
        return;
    }
    for (const summary of state.conversations) {
        const item = document.createElement("button");
        item.type = "button";
        item.className = "conversation-item";
        item.dataset.id = String(summary.user.id);
        if (state.activeConversation?.user?.id === summary.user.id) {
            item.classList.add("active");
        }
        const title = summary.user.firstName || summary.user.username || `ID ${summary.user.telegramId}`;
        const subtitle = summary.lastMessage?.body ? summary.lastMessage.body.slice(0, 80) : "Без сообщений";
        item.innerHTML = `
            <strong>${escapeHtml(title)}</strong>
            <span class="subtitle">${escapeHtml(subtitle)}</span>
        `;
        if (summary.unreadCount) {
            const badge = document.createElement("span");
            badge.className = "badge warning";
            badge.textContent = summary.unreadCount > 99 ? "99+" : String(summary.unreadCount);
            item.appendChild(badge);
        }
        el.conversationList.appendChild(item);
    }
    if (state.activeConversation) {
        const exists = state.conversations.some((c) => c.user.id === state.activeConversation.user.id);
        if (!exists) {
            state.activeConversation = null;
            renderConversationThread(null);
        }
    }
}

function renderConversationThread(conversation) {
    if (!el.conversationMessages || !el.conversationHeader || !el.conversationEmpty) {
        return;
    }
    if (!conversation) {
        el.conversationEmpty.hidden = false;
        el.conversationMessages.innerHTML = "";
        el.conversationHeader.textContent = "Выберите пользователя";
        el.conversationReplyForm?.classList.add("disabled");
        return;
    }
    state.activeConversation = conversation;
    el.conversationEmpty.hidden = true;
    const nameParts = [conversation.user.firstName, conversation.user.lastName].filter(Boolean);
    const displayName = nameParts.length ? nameParts.join(" ") : conversation.user.username || `ID ${conversation.user.telegramId}`;
    el.conversationHeader.textContent = displayName;
    el.conversationMessages.innerHTML = "";
    for (const message of conversation.messages || []) {
        const bubble = document.createElement("div");
        bubble.className = `message ${message.direction === "inbound" ? "inbound" : "outbound"}`;
        const text = document.createElement("p");
        text.textContent = message.body || "";
        const meta = document.createElement("span");
        meta.className = "meta";
        meta.textContent = formatDateTime(message.createdAt);
        bubble.append(text, meta);
        el.conversationMessages.appendChild(bubble);
    }
    el.conversationMessages.scrollTop = el.conversationMessages.scrollHeight;
    el.conversationReplyForm?.classList.remove("disabled");
}

async function loadConversation(userId, showLoader = true) {
    if (!userId) {
        return;
    }
    try {
        if (showLoader && el.conversationMessages) {
            el.conversationMessages.innerHTML = "<p class=\"hint\">Загрузка истории...</p>";
        }
        const convo = await request(`telegramConversation?id=${encodeURIComponent(userId)}`, { method: "GET" });
        renderConversationThread(convo);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    }
}

el.configForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const payload = {
        enabled: el.configEnabled.checked,
        botToken: el.botTokenInput.value.trim(),
        webhookDomain: el.webhookDomain.value.trim(),
        webhookSecret: el.webhookSecret.value.trim(),
        yooKassaShopId: el.shopId.value.trim(),
        yooKassaSecretKey: el.secretKey.value.trim(),
        successRedirectUrl: el.successUrl.value.trim(),
        failureRedirectUrl: el.failureUrl.value.trim(),
        miniAppUrl: el.miniAppUrl.value.trim(),
        downloadLinks: collectDownloadLinks(),
    };
    try {
        setLoading(true);
        await request("telegramConfig", {
            method: "POST",
            body: JSON.stringify(payload),
        });
        showToast("Настройки сохранены");
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
});

el.addDownloadLink.addEventListener("click", () => addDownloadLinkRow());

el.createTariff.addEventListener("click", () => {
    state.selectedTariffId = null;
    for (const row of el.tariffTableBody.querySelectorAll("tr")) {
        row.classList.remove("selected");
    }
    resetTariffForm();
    renderButtons();
});

el.refreshTariffs.addEventListener("click", () => loadState(true));

el.tariffTableBody.addEventListener("click", (event) => {
    const button = event.target.closest("button[data-action]");
    if (button) {
        event.stopPropagation();
        const id = Number(button.dataset.id);
        if (button.dataset.action === "edit") {
            const tariff = state.tariffs.find((t) => t.id === id);
            if (tariff) {
                state.selectedTariffId = id;
                populateTariffForm(tariff);
                renderTariffs();
                renderButtons();
            }
        } else if (button.dataset.action === "delete") {
            confirmDeleteTariff(id);
        }
        return;
    }
    const row = event.target.closest("tr[data-id]");
    if (row) {
        selectTariff(row.dataset.id);
        populateTariffForm(state.tariffs.find((t) => t.id === Number(row.dataset.id)) || {});
    }
});

async function confirmDeleteTariff(id) {
    if (!confirm("Удалить тариф и связанные кнопки?")) {
        return;
    }
    try {
        setLoading(true);
        await request("telegramTariffDelete", {
            method: "POST",
            body: JSON.stringify({ id }),
        });
        showToast("Тариф удалён");
        if (state.selectedTariffId === id) {
            state.selectedTariffId = null;
        }
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
}

el.tariffForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const priceValue = parseFloat(el.tariffPrice.value);
    if (Number.isNaN(priceValue) || priceValue <= 0) {
        showToast("Укажите корректную цену", true);
        return;
    }
    const payload = {
        id: el.tariffId.value ? Number(el.tariffId.value) : 0,
        title: el.tariffTitle.value.trim(),
        description: el.tariffDescription.value.trim(),
        priceMinor: Math.round(priceValue * 100),
        currency: el.tariffCurrency.value.trim().toUpperCase(),
        durationDays: el.tariffDuration.value ? Number(el.tariffDuration.value) : 0,
        sortOrder: el.tariffSort.value ? Number(el.tariffSort.value) : 0,
        active: el.tariffActive.checked,
    };
    try {
        setLoading(true);
        const obj = await request("telegramTariff", {
            method: "POST",
            body: JSON.stringify(payload),
        });
        showToast("Тариф сохранён");
        state.selectedTariffId = obj?.id || payload.id || null;
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
});

el.buttonTableBody.addEventListener("click", (event) => {
    const button = event.target.closest("button[data-action]");
    if (!button) {
        return;
    }
    event.stopPropagation();
    const id = Number(button.dataset.id);
    const tariff = state.tariffs.find((t) => t.id === state.selectedTariffId);
    if (!tariff) {
        return;
    }
    const targetButton = tariff.buttons?.find((b) => b.id === id);
    if (button.dataset.action === "edit") {
        if (targetButton) {
            openButtonForm(targetButton);
        }
    } else if (button.dataset.action === "delete") {
        confirmDeleteButton(id);
    }
});

async function confirmDeleteButton(id) {
    if (!confirm("Удалить кнопку?")) {
        return;
    }
    try {
        setLoading(true);
        await request("telegramButtonDelete", {
            method: "POST",
            body: JSON.stringify({ id }),
        });
        showToast("Кнопка удалена");
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
}

el.buttonForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    if (!state.selectedTariffId) {
        showToast("Выберите тариф", true);
        return;
    }
    const payload = {
        id: el.buttonId.value ? Number(el.buttonId.value) : 0,
        tariffId: el.buttonTariffId.value ? Number(el.buttonTariffId.value) : state.selectedTariffId,
        label: el.buttonLabel.value.trim(),
        action: el.buttonAction.value.trim(),
        payload: el.buttonPayload.value.trim(),
        sortOrder: el.buttonSort.value ? Number(el.buttonSort.value) : 0,
    };
    if (!payload.label || !payload.action) {
        showToast("Укажите подпись и действие", true);
        return;
    }
    try {
        setLoading(true);
        await request("telegramButton", {
            method: "POST",
            body: JSON.stringify(payload),
        });
        showToast("Кнопка сохранена");
        closeButtonForm();
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
});

el.cancelButtonEdit.addEventListener("click", () => closeButtonForm());

el.addButton.addEventListener("click", () => openButtonForm());

async function saveBroadcast(payload) {
    try {
        setLoading(true);
        const obj = await request("telegramBroadcast", {
            method: "POST",
            body: JSON.stringify(payload),
        });
        showToast("Рассылка сохранена");
        state.selectedBroadcastId = obj?.id || payload.id || null;
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
}

async function sendBroadcast(id) {
    if (!id) {
        showToast("Сначала сохраните рассылку", true);
        return;
    }
    if (!confirm("Отправить рассылку выбранным пользователям?")) {
        return;
    }
    try {
        setLoading(true);
        await request("telegramBroadcastSend", {
            method: "POST",
            body: JSON.stringify({ id }),
        });
        showToast("Рассылка запущена");
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
}

async function deleteBroadcast(id) {
    if (!id) {
        return;
    }
    if (!confirm("Удалить черновик рассылки?")) {
        return;
    }
    try {
        setLoading(true);
        await request("telegramBroadcastDelete", {
            method: "POST",
            body: JSON.stringify({ id }),
        });
        showToast("Рассылка удалена");
        state.selectedBroadcastId = null;
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
}

async function editSentBroadcast(id) {
    if (!id) {
        return;
    }
    const body = el.broadcastBody?.value.trim();
    if (!body) {
        showToast("Введите текст сообщения", true);
        return;
    }
    try {
        setLoading(true);
        await request("telegramBroadcastEdit", {
            method: "POST",
            body: JSON.stringify({ broadcastId: id, body }),
        });
        showToast("Сообщения обновлены");
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
}

async function savePromo(payload) {
    try {
        setLoading(true);
        await request("telegramPromo", {
            method: "POST",
            body: JSON.stringify(payload),
        });
        showToast("Промокод сохранён");
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
}

async function deletePromo(id) {
    if (!id) {
        return;
    }
    if (!confirm("Удалить промокод?")) {
        return;
    }
    try {
        setLoading(true);
        await request("telegramPromoDelete", {
            method: "POST",
            body: JSON.stringify({ id }),
        });
        showToast("Промокод удалён");
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
}

async function replyConversation(id, text) {
    if (!id) {
        showToast("Выберите пользователя", true);
        return;
    }
    if (!text.trim()) {
        showToast("Введите сообщение", true);
        return;
    }
    try {
        setLoading(true);
        const convo = await request("telegramConversationReply", {
            method: "POST",
            body: JSON.stringify({ id, text }),
        });
        showToast("Сообщение отправлено");
        renderConversationThread(convo);
        el.conversationReply.value = "";
        await loadState(true);
    } catch (error) {
        console.error(error);
        showToast(error.message, true);
    } finally {
        setLoading(false);
    }
}

if (el.broadcastForm) {
    el.broadcastForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const payload = collectBroadcastPayload();
        if (!payload.title || !payload.body) {
            showToast("Укажите название и текст рассылки", true);
            return;
        }
        if (!payload.audience.allUsers && !payload.audience.tariffIds.length && !payload.audience.includeNeverSubscribed && !payload.audience.includeExpired) {
            showToast("Выберите хотя бы один сегмент аудитории", true);
            return;
        }
        await saveBroadcast(payload);
    });
}

if (el.broadcastCreate) {
    el.broadcastCreate.addEventListener("click", () => {
        state.selectedBroadcastId = null;
        for (const row of el.broadcastTableBody.querySelectorAll("tr")) {
            row.classList.remove("selected");
        }
        resetBroadcastForm();
    });
}

if (el.broadcastRefresh) {
    el.broadcastRefresh.addEventListener("click", () => loadState(true));
}

if (el.broadcastTableBody) {
    el.broadcastTableBody.addEventListener("click", (event) => {
        const button = event.target.closest("button[data-action]");
        if (!button) {
            const row = event.target.closest("tr[data-id]");
            if (row) {
                const broadcast = state.broadcasts.find((b) => b.id === Number(row.dataset.id));
                if (broadcast) {
                    state.selectedBroadcastId = broadcast.id;
                    populateBroadcastForm(broadcast);
                    renderBroadcasts();
                }
            }
            return;
        }
        event.stopPropagation();
        const id = Number(button.dataset.id);
        if (button.dataset.action === "edit") {
            const broadcast = state.broadcasts.find((b) => b.id === id);
            if (broadcast) {
                state.selectedBroadcastId = broadcast.id;
                populateBroadcastForm(broadcast);
                renderBroadcasts();
            }
        } else if (button.dataset.action === "send") {
            sendBroadcast(id);
        } else if (button.dataset.action === "delete") {
            deleteBroadcast(id);
        } else if (button.dataset.action === "edit-sent") {
            editSentBroadcast(id);
        }
    });
}

if (el.broadcastSend) {
    el.broadcastSend.addEventListener("click", () => {
        const id = Number(el.broadcastId?.value || 0);
        sendBroadcast(id);
    });
}

if (el.broadcastEditSent) {
    el.broadcastEditSent.addEventListener("click", () => {
        const id = Number(el.broadcastId?.value || 0);
        editSentBroadcast(id);
    });
}

if (el.broadcastAllUsers) {
    el.broadcastAllUsers.addEventListener("change", () => updateBroadcastAudienceControls());
}

if (el.broadcastTariffOptions) {
    el.broadcastTariffOptions.addEventListener("change", () => {
        if (el.broadcastAllUsers && el.broadcastAllUsers.checked) {
            updateBroadcastAudienceControls();
        }
    });
}

if (el.promoForm) {
    el.promoForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const payload = collectPromoPayload();
        if (!payload.code) {
            showToast("Укажите код", true);
            return;
        }
        await savePromo(payload);
    });
}

if (el.promoTableBody) {
    el.promoTableBody.addEventListener("click", (event) => {
        const button = event.target.closest("button[data-action]");
        if (!button) {
            return;
        }
        event.stopPropagation();
        const id = Number(button.dataset.id);
        if (button.dataset.action === "edit") {
            const promo = state.promos.find((p) => p.id === id);
            if (promo) {
                populatePromoForm(promo);
            }
        } else if (button.dataset.action === "delete") {
            deletePromo(id);
        }
    });
}

if (el.promoNoExpiry) {
    el.promoNoExpiry.addEventListener("change", () => {
        el.promoExpiry.disabled = el.promoNoExpiry.checked;
        if (el.promoNoExpiry.checked) {
            el.promoExpiry.value = "";
        }
    });
}

if (el.refreshConversations) {
    el.refreshConversations.addEventListener("click", () => loadState(true));
}

if (el.conversationList) {
    el.conversationList.addEventListener("click", (event) => {
        const button = event.target.closest("button.conversation-item");
        if (!button) {
            return;
        }
        const id = Number(button.dataset.id);
        loadConversation(id);
    });
}

if (el.conversationReplyForm) {
    el.conversationReplyForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const id = state.activeConversation?.user?.id;
        await replyConversation(id, el.conversationReply.value);
    });
}

window.addEventListener("pageshow", () => {
    showToast("", false);
});

loadState();
