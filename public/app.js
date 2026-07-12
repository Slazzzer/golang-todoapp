'use strict';

/* =========================================================================
 * Конфигурация адреса API
 * -------------------------------------------------------------------------
 * Приоритет: ?api=... в URL  →  сохранённое значение в localStorage  →
 * тот же origin (если страница открыта по http/https)  →  http://localhost:5050
 * ========================================================================= */
const DEFAULT_API = 'http://localhost:5050';
const API_STORE_KEY = 'todoapp_api_base';
const ACTING_USER_KEY = 'todoapp_acting_user_id';
const ADMIN_DISPLAY_NAME = 'Admin';
const ACTING_ROLE_KEY = 'todoapp_acting_role';
const ADMIN_SESSION_KEY = 'todoapp_admin_session';

function isAdmin() {
    return localStorage.getItem(ACTING_ROLE_KEY) === 'admin' && !!getAdminSession();
}

function getAdminSession() {
    return localStorage.getItem(ADMIN_SESSION_KEY);
}

function setAdminSession(session) {
    localStorage.setItem(ADMIN_SESSION_KEY, session);
    localStorage.setItem(ACTING_ROLE_KEY, 'admin');
    localStorage.removeItem(ACTING_USER_KEY);
}

function clearAdminSession() {
    localStorage.removeItem(ADMIN_SESSION_KEY);
    if (localStorage.getItem(ACTING_ROLE_KEY) === 'admin') {
        localStorage.removeItem(ACTING_ROLE_KEY);
    }
}

function clearActingSession() {
    clearAdminSession();
    clearActingUserId();
    localStorage.removeItem(ACTING_ROLE_KEY);
}

function getActingUserId() {
    const raw = localStorage.getItem(ACTING_USER_KEY);
    if (!raw) return null;
    const id = parseInt(raw, 10);
    return Number.isFinite(id) && id > 0 ? id : null;
}

function setActingUserId(id) {
    localStorage.setItem(ACTING_USER_KEY, String(id));
    localStorage.setItem(ACTING_ROLE_KEY, 'user');
    clearAdminSession();
}

function clearActingUserId() {
    localStorage.removeItem(ACTING_USER_KEY);
}

function resolveApiBase() {
    const fromQuery = new URLSearchParams(location.search).get('api');
    if (fromQuery) {
        localStorage.setItem(API_STORE_KEY, fromQuery);
        return fromQuery.replace(/\/+$/, '');
    }
    const stored = localStorage.getItem(API_STORE_KEY);
    if (stored) return stored.replace(/\/+$/, '');
    // Фронт и API на одном origin только если открыто с todoapp (:5050).
    if (location.protocol === 'http:' || location.protocol === 'https:') {
        const port = location.port;
        if (port === '5050' || port === '') return location.origin;
    }
    return DEFAULT_API;
}

let API_BASE = resolveApiBase();
const apiV1 = () => `${API_BASE}/api/v1`;

/* =========================================================================
 * HTTP-хелперы
 * ========================================================================= */
function humanizeError(text, status) {
    if (!text) return `Ошибка сервера (${status})`;

    try {
        const data = JSON.parse(text);
        if (data && data.message) return String(data.message);
    } catch { /* не JSON */ }

    // HTML-ответ (404 от статик-сервера и т.п.) — не показываем разметку
    if (/<!DOCTYPE|<html[\s>]/i.test(text)) {
        const pre = text.match(/<pre[^>]*>([\s\S]*?)<\/pre>/i);
        if (pre) {
            const msg = pre[1].replace(/<[^>]+>/g, '').trim();
            if (/Cannot (GET|POST|PATCH|DELETE|PUT)/i.test(msg)) {
                return 'Сервер API недоступен. Запустите backend (порт 5050) или укажите ?api=http://localhost:5050';
            }
            return msg.length > 140 ? msg.slice(0, 140) + '…' : msg;
        }
        if (status === 404) return 'Сервер API недоступен. Проверьте, что backend запущен.';
        return `Ошибка сервера (${status})`;
    }

    const clean = text.replace(/\s+/g, ' ').trim();
    return clean.length > 140 ? clean.slice(0, 140) + '…' : clean;
}

async function request(method, path, body, options = {}) {
    const headers = {};
    if (!options.public) {
        if (isAdmin()) {
            headers['X-Admin-Session'] = getAdminSession();
        } else {
            const uid = getActingUserId();
            if (uid) headers['X-User-ID'] = String(uid);
        }
    }
    const reqOptions = { method, headers, cache: 'no-store' };
    if (body !== undefined) {
        reqOptions.headers['Content-Type'] = 'application/json';
        reqOptions.body = JSON.stringify(body);
    }

    let response;
    try {
        response = await fetch(`${apiV1()}${path}`, reqOptions);
    } catch {
        throw new Error('Нет связи с сервером. Проверьте, что backend запущен на порту 5050.');
    }

    if (response.status === 204) return null;

    const text = await response.text();
    if (!response.ok) {
        let body = null;
        try { body = JSON.parse(text); } catch { /* не JSON */ }
        const error = new Error(humanizeError(text, response.status));
        error.status = response.status;
        error.code = body && body.code ? body.code : '';
        throw error;
    }

    if (!text) return null;
    try { return JSON.parse(text); }
    catch { throw new Error('Сервер вернул некорректный ответ.'); }
}

const api = {
    loginAdmin: (credentials) => request('POST', '/auth/login', credentials, { public: true }),
    listUsersPublic: (params) => request('GET', `/users${query(params)}`, undefined, { public: true }),
    createUserPublic: (u) => request('POST', '/users', u, { public: true }),
    listUsers: (params) => request('GET', `/users${query(params)}`),
    createUser: (u) => request('POST', '/users', u),
    patchUser: (id, u) => request('PATCH', `/users/${id}`, u),
    deleteUser: (id) => request('DELETE', `/users/${id}`),
    listTasks: (params) => request('GET', `/tasks${query(params)}`),
    createTask: (t) => request('POST', '/tasks', t),
    patchTask: (id, t) => request('PATCH', `/tasks/${id}`, t),
    deleteTask: (id) => request('DELETE', `/tasks/${id}`),
    statistics: (params) => request('GET', `/statistics${query(params)}`),
};

function query(params) {
    const q = new URLSearchParams();
    Object.entries(params || {}).forEach(([k, v]) => {
        if (v !== undefined && v !== null && v !== '') q.append(k, v);
    });
    const s = q.toString();
    return s ? `?${s}` : '';
}

/* =========================================================================
 * Утилиты
 * ========================================================================= */
const $ = (sel, root = document) => root.querySelector(sel);
const $$ = (sel, root = document) => Array.from(root.querySelectorAll(sel));

function el(tag, attrs = {}, children = []) {
    const node = document.createElement(tag);
    for (const [k, v] of Object.entries(attrs)) {
        if (k === 'class') node.className = v;
        else if (k === 'html') node.innerHTML = v;
        else if (k === 'text') node.textContent = v;
        else if (k.startsWith('on') && typeof v === 'function') node.addEventListener(k.slice(2), v);
        else if (v !== null && v !== undefined) node.setAttribute(k, v);
    }
    (Array.isArray(children) ? children : [children]).forEach((c) => {
        if (c == null) return;
        node.appendChild(typeof c === 'string' ? document.createTextNode(c) : c);
    });
    return node;
}

function icon(id, cls = 'ic') {
    return `<svg class="${cls}"><use href="#${id}"/></svg>`;
}

function escapeHtml(str) {
    return String(str).replace(/[&<>"']/g, (c) => (
        { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]
    ));
}

function formatDate(iso) {
    if (!iso) return '—';
    const d = new Date(iso);
    if (isNaN(d)) return '—';
    return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short', year: 'numeric' });
}

function relativeAge(iso) {
    const then = new Date(iso).getTime();
    if (isNaN(then)) return '';
    const sec = Math.max(0, Math.floor((Date.now() - then) / 1000));
    if (sec < 60) return `${sec} с`;
    const min = Math.floor(sec / 60);
    if (min < 60) return `${min} мин`;
    const hr = Math.floor(min / 60);
    if (hr < 24) return `${hr} ч`;
    const day = Math.floor(hr / 24);
    return `${day} дн`;
}

function formatDuration(goDuration) {
    if (!goDuration) return '—';
    // Go-длительность вида "26h3m4.5s" → человекочитаемо
    const m = /(?:(\d+)h)?(?:(\d+)m)?(?:([\d.]+)s)?/.exec(goDuration);
    if (!m) return goDuration;
    const h = parseInt(m[1] || '0', 10);
    const min = parseInt(m[2] || '0', 10);
    const s = Math.round(parseFloat(m[3] || '0'));
    const parts = [];
    if (h) parts.push(`${h} ч`);
    if (min) parts.push(`${min} мин`);
    if (s && !h) parts.push(`${s} с`);
    return parts.length ? parts.join(' ') : '0 с';
}

const AVATAR_COLORS = [
    'linear-gradient(135deg,#f97316,#ea580c)',
    'linear-gradient(135deg,#a855f7,#d946ef)',
    'linear-gradient(135deg,#3b82f6,#2563eb)',
    'linear-gradient(135deg,#10b981,#059669)',
    'linear-gradient(135deg,#ec4899,#db2777)',
    'linear-gradient(135deg,#f59e0b,#d97706)',
    'linear-gradient(135deg,#06b6d4,#0891b2)',
];
function avatarColor(id) { return AVATAR_COLORS[id % AVATAR_COLORS.length]; }
function initials(name) {
    const parts = String(name).trim().split(/\s+/).slice(0, 2);
    return parts.map((p) => p[0] ? p[0].toUpperCase() : '').join('') || '?';
}

/* =========================================================================
 * Тосты
 * ========================================================================= */
function toast(title, msg = '', kind = '') {
    const node = el('div', { class: `toast ${kind}` }, [
        el('div', { class: 'toast-title', text: title }),
        msg ? el('div', { class: 'toast-msg', text: msg }) : null,
    ]);
    $('#toasts').appendChild(node);
    setTimeout(() => {
        node.style.transition = 'opacity .2s, transform .2s';
        node.style.opacity = '0';
        node.style.transform = 'translateX(20px)';
        setTimeout(() => node.remove(), 220);
    }, 3500);
}
const toastOk = (t, m) => toast(t, m, 'ok');
const toastErr = (t, m) => toast(t, m, 'err');

/* =========================================================================
 * Состояние
 * ========================================================================= */
const PAGE_SIZE_OPTIONS = [10, 20, 50, 100];

const state = {
    users: [],
    actingUser: null,
    editingUserId: null,
    editingTaskId: null,
    tasksPage: 1,
    usersPage: 1,
    taskPageSize: 20,
    userPageSize: 20,
    identityRegisterMode: false,
    identityAdminMode: false,
    identityCanCancel: false,
    identityListLoadId: 0,
};

async function loadActingUser() {
    if (isAdmin()) {
        state.actingUser = null;
        await loadUsersForAdminSelects();
        updateRoleUI();
        return null;
    }

    const id = getActingUserId();
    if (!id) {
        state.actingUser = null;
        state.users = [];
        return null;
    }

    const page = await api.listUsers({ limit: 1, offset: 0 });
    const user = (page.items || [])[0] || null;
    state.actingUser = user;
    state.users = user ? [user] : [];
    updateRoleUI();
    return user;
}

function updateBrandBadge() {
    const badge = $('#brandBadge');
    if (!badge) return;
    if (isAdmin()) {
        badge.textContent = ADMIN_DISPLAY_NAME;
        return;
    }
    badge.textContent = state.actingUser ? state.actingUser.full_name : '—';
}

function updateRoleUI() {
    const roleValue = $('#navRoleValue');
    const usersNavLabel = $('.nav-item[data-view="users"] span');
    const usersTitle = $('#usersPageTitle');
    const statFilterWrap = $('#statUserFilterWrap');
    const newTaskBtn = $('#newTaskBtn');

    if (isAdmin()) {
        if (roleValue) roleValue.textContent = ADMIN_DISPLAY_NAME;
        if (usersNavLabel) usersNavLabel.textContent = 'Пользователи';
        if (usersTitle) usersTitle.textContent = 'Пользователи';
        if (statFilterWrap) statFilterWrap.classList.remove('is-hidden');
        if (newTaskBtn) newTaskBtn.classList.add('is-hidden');
    } else {
        if (roleValue) roleValue.textContent = 'Пользователь';
        if (usersNavLabel) usersNavLabel.textContent = 'Мой профиль';
        if (usersTitle) usersTitle.textContent = 'Мой профиль';
        if (statFilterWrap) statFilterWrap.classList.add('is-hidden');
        if (newTaskBtn) newTaskBtn.classList.remove('is-hidden');
    }

    updateBrandBadge();
}

async function loadUsersForAdminSelects() {
    const all = [];
    let offset = 0;
    const limit = 100;

    while (true) {
        const page = await api.listUsers({ limit, offset });
        const items = page.items || [];
        if (!items.length) break;
        all.push(...items);
        const total = Number(page.total);
        if (Number.isFinite(total) && total > 0 && all.length >= total) break;
        offset += items.length;
    }

    state.users = all;
    const userOptions = all
        .map((u) => `<option value="${u.id}" data-user-id="${u.id}">${escapeHtml(u.full_name)}</option>`)
        .join('');
    const statFilter = $('#statUserFilter');
    if (statFilter) {
        statFilter.innerHTML = `<option value="">Все пользователи</option>${userOptions}`;
        refreshCustomSelect(statFilter);
    }
    return all;
}

async function refreshUsersCatalog() {
    await loadActingUser();
}

function userById(id) { return state.users.find((u) => u.id === id) || (state.actingUser && state.actingUser.id === id ? state.actingUser : null); }

function userFromOption(opt) {
    if (!opt || !opt.dataset.userId) return null;
    return userById(parseInt(opt.dataset.userId, 10));
}

function userOptionNode(u) {
    const node = el('span', { class: 'cselect-user' });
    node.innerHTML = `
        <span class="cselect-user-avatar" style="background:${avatarColor(u.id)}">${escapeHtml(initials(u.full_name))}</span>
        <span class="cselect-user-text">
            <span class="cselect-user-name">${escapeHtml(u.full_name)}</span>
            <span class="cselect-user-id">ID ${u.id}</span>
        </span>`;
    return node;
}

/* =========================================================================
 * Выбор пользователя (acting user)
 * ========================================================================= */
async function loadAllUsersPublic() {
    const all = [];
    let offset = 0;
    const limit = 100;

    while (true) {
        const page = await api.listUsersPublic({ limit, offset });
        const items = page.items || [];
        if (!items.length) break;

        all.push(...items);
        if (items.length < limit) break;

        offset += items.length;
    }

    return all;
}

function syncIdentityFooter() {
    const inSubMode = state.identityRegisterMode || state.identityAdminMode;

    $('#identityListActions').classList.toggle('is-hidden', inSubMode);
    $('#identityBackBtn').classList.toggle('is-hidden', !inSubMode);
    $('#identityCancelBtn').classList.toggle('is-hidden', inSubMode || !state.identityCanCancel);
    $('#identityCloseBtn').classList.toggle('is-hidden', !inSubMode && !state.identityCanCancel);
}

function setIdentityRegisterMode(on) {
    state.identityRegisterMode = on;
    if (on) state.identityAdminMode = false;
    $('#identityRegister').classList.toggle('is-hidden', !on);
    $('#identityAdmin').classList.toggle('is-hidden', true);
    $('#identityList').classList.toggle('is-hidden', on);
    $('.identity-hint').classList.toggle('is-hidden', on);
    $('#identitySubmitRegister').classList.toggle('is-hidden', !on);
    $('#identitySubmitAdmin').classList.add('is-hidden');
    hideError('identityFormError');
    hideError('identityAdminError');
    syncIdentityFooter();
}

function setIdentityAdminMode(on) {
    state.identityAdminMode = on;
    if (on) state.identityRegisterMode = false;
    $('#identityAdmin').classList.toggle('is-hidden', !on);
    $('#identityRegister').classList.toggle('is-hidden', true);
    $('#identityList').classList.toggle('is-hidden', on);
    $('.identity-hint').classList.toggle('is-hidden', on);
    $('#identitySubmitAdmin').classList.toggle('is-hidden', !on);
    $('#identitySubmitRegister').classList.add('is-hidden');
    hideError('identityFormError');
    hideError('identityAdminError');
    syncIdentityFooter();
}

function backToIdentityList({ refreshList = false } = {}) {
    setIdentityRegisterMode(false);
    setIdentityAdminMode(false);
    $('#identityFullName').value = '';
    $('#identityPhone').value = '';
    $('#identityAdminLogin').value = '';
    $('#identityAdminPassword').value = '';
    hideError('identityFormError');
    hideError('identityAdminError');
    if (refreshList && !$('#identityModal').classList.contains('is-hidden')) {
        renderIdentityList();
    }
}

async function renderIdentityList() {
    const loadId = ++state.identityListLoadId;
    const list = $('#identityList');
    list.innerHTML = '';
    list.appendChild(el('div', { class: 'empty', html: `${icon('i-clock', 'ic')}<p>Загрузка…</p>` }));

    let users;
    try {
        users = await loadAllUsersPublic();
    } catch (err) {
        if (loadId !== state.identityListLoadId) return;
        list.innerHTML = '';
        list.appendChild(emptyState('i-users', 'Не удалось загрузить пользователей', err.message));
        return;
    }

    if (loadId !== state.identityListLoadId) return;

    list.innerHTML = '';
    if (!users.length) {
        list.appendChild(emptyState('i-users', 'Пользователей пока нет', 'Создайте первого пользователя'));
        setIdentityRegisterMode(true);
        return;
    }

    users.forEach((u) => {
        const btn = el('button', { class: 'identity-item', type: 'button' });
        btn.innerHTML = `
            <span class="user-avatar" style="background:${avatarColor(u.id)}">${escapeHtml(initials(u.full_name))}</span>
            <span class="identity-item-text">
                <span class="identity-item-name">${escapeHtml(u.full_name)}</span>
                <span class="identity-item-meta">ID ${u.id}${u.phone_number ? ` · ${escapeHtml(u.phone_number)}` : ''}</span>
            </span>`;
        btn.addEventListener('click', () => selectActingUser(u.id));
        list.appendChild(btn);
    });
}

function openIdentityModal({ canCancel = false } = {}) {
    state.identityCanCancel = canCancel;
    backToIdentityList();
    $('#identityModal').classList.remove('is-hidden');
    renderIdentityList();
}

function closeIdentityModal() {
    const restoreFocus = state.identityCanCancel;
    backToIdentityList();
    state.identityCanCancel = false;
    $('#identityModal').classList.add('is-hidden');
    if (restoreFocus) {
        const btn = $('#switchUserBtn');
        if (btn) btn.focus({ preventScroll: true });
    }
}

function handleIdentityDismiss() {
    if (state.identityRegisterMode || state.identityAdminMode) {
        backToIdentityList({ refreshList: true });
        const focusTarget = $('#identityToggleRegister') || $('#identityList .identity-item');
        if (focusTarget) focusTarget.focus({ preventScroll: true });
        return;
    }
    if (state.identityCanCancel) closeIdentityModal();
}

async function selectActingUser(id, options = {}) {
    setActingUserId(id);
    closeIdentityModal();
    try {
        const user = await loadActingUser();
        if (!options.silentWelcome && user) {
            toastOk(`Добро пожаловать, ${user.full_name}`);
        }
        switchView('tasks');
    } catch (err) {
        clearActingSession();
        toastErr('Не удалось войти', err.message);
        openIdentityModal();
    }
}

async function submitIdentityAdmin() {
    const login = $('#identityAdminLogin').value.trim();
    const password = $('#identityAdminPassword').value;
    hideError('identityAdminError');

    if (!login || !password) {
        return showError('identityAdminError', 'Введите логин и пароль.');
    }

    const btn = $('#identitySubmitAdmin');
    btn.disabled = true;
    try {
        const result = await api.loginAdmin({ login, password });
        setAdminSession(result.session);
        closeIdentityModal();
        await loadActingUser();
        toastOk(`Добро пожаловать, ${ADMIN_DISPLAY_NAME}`);
        switchView('users');
    } catch (err) {
        showError('identityAdminError', err.message);
    } finally {
        btn.disabled = false;
    }
}

async function submitIdentityRegister() {
    const fullName = $('#identityFullName').value.trim();
    const phone = $('#identityPhone').value.trim();
    hideError('identityFormError');

    if (fullName.length < 3) return showError('identityFormError', 'Полное имя должно содержать минимум 3 символа.');
    if (phone && !phone.startsWith('+')) return showError('identityFormError', 'Телефон должен начинаться с «+».');
    if (phone && (phone.length < 10 || phone.length > 15)) {
        return showError('identityFormError', 'Телефон должен быть от 10 до 15 символов.');
    }

    const btn = $('#identitySubmitRegister');
    btn.disabled = true;
    try {
        const body = { full_name: fullName };
        if (phone) body.phone_number = phone;
        const created = await api.createUserPublic(body);
        await selectActingUser(created.id, { silentWelcome: true });
        toastOk(`Добро пожаловать, ${fullName}`);
    } catch (err) {
        showError('identityFormError', err.message);
    } finally {
        btn.disabled = false;
    }
}

function switchActingUser() {
    openIdentityModal({ canCancel: true });
}

/* =========================================================================
 * Кастомные выпадающие списки
 * ========================================================================= */
const customSelectMap = new WeakMap();

function closeAllCustomSelects() {
    $$('.cselect-menu').forEach((m) => m.classList.add('is-hidden'));
    $$('.cselect').forEach((c) => c.classList.remove('is-open'));
}

function enhanceSelect(nativeSel) {
    if (!nativeSel || nativeSel.dataset.cselectDone) return;
    nativeSel.dataset.cselectDone = '1';
    nativeSel.classList.add('cselect-native');

    const wrap = el('div', { class: 'cselect' });
    nativeSel.parentNode.insertBefore(wrap, nativeSel);
    wrap.appendChild(nativeSel);

    const btn = el('button', { class: 'cselect-btn', type: 'button' });
    const label = el('span', { class: 'cselect-label' });
    const chevron = el('span', { class: 'cselect-chevron', 'aria-hidden': 'true' });
    btn.appendChild(label);
    btn.appendChild(chevron);

    const menu = el('div', { class: 'cselect-menu is-hidden' });
    wrap.appendChild(btn);
    wrap.appendChild(menu);

    function syncLabel() {
        const opt = nativeSel.options[nativeSel.selectedIndex];
        label.innerHTML = '';
        if (!opt) return;
        const u = userFromOption(opt);
        if (u) label.appendChild(userOptionNode(u));
        else label.textContent = opt.textContent;
    }

    function buildMenu() {
        menu.innerHTML = '';
        Array.from(nativeSel.options).forEach((opt, i) => {
            const u = userFromOption(opt);
            const item = el('button', {
                class: 'cselect-option' + (opt.selected ? ' is-selected' : ''),
                type: 'button',
                onclick: () => selectOption(i),
            });
            if (u) item.appendChild(userOptionNode(u));
            else item.textContent = opt.textContent;
            menu.appendChild(item);
        });
    }

    function selectOption(index) {
        nativeSel.selectedIndex = index;
        syncLabel();
        buildMenu();
        closeAllCustomSelects();
        nativeSel.dispatchEvent(new Event('change', { bubbles: true }));
    }

    btn.addEventListener('click', (e) => {
        e.stopPropagation();
        if (nativeSel.disabled) return;
        const isOpen = wrap.classList.contains('is-open');
        closeAllCustomSelects();
        if (!isOpen) {
            buildMenu();
            menu.classList.remove('is-hidden');
            wrap.classList.add('is-open');
        }
    });

    customSelectMap.set(nativeSel, { buildMenu, syncLabel });
    buildMenu();
    syncLabel();
}

function refreshCustomSelect(nativeSel) {
    const meta = customSelectMap.get(nativeSel);
    const wrap = nativeSel.closest('.cselect');
    if (wrap) wrap.classList.toggle('is-disabled', nativeSel.disabled);

    if (meta) {
        meta.buildMenu();
        meta.syncLabel();
    } else {
        enhanceSelect(nativeSel);
    }
}

function setSelectValue(nativeSel, value) {
    nativeSel.value = value;
    refreshCustomSelect(nativeSel);
}

function initCustomSelects() {
    $$('select').forEach(enhanceSelect);
    document.addEventListener('click', closeAllCustomSelects);
}

function setContentEmpty(areaId, isEmpty) {
    const area = $(`#${areaId}`);
    if (area) area.classList.toggle('is-empty', isEmpty);
}

function totalPages(total, pageSize) {
    return Math.max(1, Math.ceil(total / pageSize));
}

function pageNumberWindow(current, total) {
    if (total <= 7) {
        return Array.from({ length: total }, (_, i) => i + 1);
    }

    const nums = new Set([1, total, current - 1, current, current + 1]);
    const sorted = [...nums].filter((n) => n >= 1 && n <= total).sort((a, b) => a - b);
    const out = [];
    let prev = 0;

    sorted.forEach((n) => {
        if (prev && n - prev > 1) out.push('…');
        out.push(n);
        prev = n;
    });

    return out;
}

function renderPaginationBar(container, { page, total, pageSize, selectId, onPage, onPageSize }) {
    container.innerHTML = '';

    if (!total) {
        container.classList.add('is-hidden');
        return;
    }

    const pages = totalPages(total, pageSize);
    const safePage = Math.min(Math.max(1, page), pages);
    const start = (safePage - 1) * pageSize + 1;
    const end = Math.min(safePage * pageSize, total);

    container.classList.remove('is-hidden');

    const info = el('div', { class: 'pagination-info' });
    info.innerHTML = `Страница <b>${safePage}</b> из <b>${pages}</b> · записи ${start}–${end} из ${total}`;

    const nav = el('div', { class: 'pagination-nav' });

    const prev = el('button', { class: 'pagination-btn', type: 'button', title: 'Предыдущая страница' });
    prev.textContent = '‹';
    prev.disabled = safePage <= 1;
    prev.addEventListener('click', () => onPage(safePage - 1));
    nav.appendChild(prev);

    pageNumberWindow(safePage, pages).forEach((item) => {
        if (item === '…') {
            nav.appendChild(el('span', { class: 'pagination-btn pagination-btn--ghost', text: '…' }));
            return;
        }
        const btn = el('button', {
            class: `pagination-btn${item === safePage ? ' is-active' : ''}`,
            type: 'button',
            text: String(item),
        });
        btn.addEventListener('click', () => { if (item !== safePage) onPage(item); });
        nav.appendChild(btn);
    });

    const next = el('button', { class: 'pagination-btn', type: 'button', title: 'Следующая страница' });
    next.textContent = '›';
    next.disabled = safePage >= pages;
    next.addEventListener('click', () => onPage(safePage + 1));
    nav.appendChild(next);

    const sizeWrap = el('div', { class: 'pagination-size' });
    sizeWrap.appendChild(el('label', { text: 'На странице' }));
    const sizeSel = el('select', { id: selectId });
    PAGE_SIZE_OPTIONS.forEach((n) => {
        const opt = document.createElement('option');
        opt.value = String(n);
        opt.textContent = String(n);
        if (n === pageSize) opt.selected = true;
        sizeSel.appendChild(opt);
    });
    sizeSel.addEventListener('change', () => onPageSize(parseInt(sizeSel.value, 10) || 20));
    sizeWrap.appendChild(sizeSel);

    container.append(info, nav, sizeWrap);
    enhanceSelect(sizeSel);
}

/* =========================================================================
 * Навигация между разделами
 * ========================================================================= */
const views = { tasks: '#view-tasks', users: '#view-users', stats: '#view-stats' };

function switchView(name) {
    Object.entries(views).forEach(([key, sel]) => {
        $(sel).classList.toggle('is-hidden', key !== name);
    });
    $$('.nav-item').forEach((b) => b.classList.toggle('is-active', b.dataset.view === name));

    if (name === 'tasks') renderTasks();
    if (name === 'users') renderUsers();
    if (name === 'stats') renderStats();
}

/* =========================================================================
 * Раздел: ЗАДАЧИ
 * ========================================================================= */
async function renderTasks() {
    const list = $('#taskList');
    const pagination = $('#taskPagination');
    list.innerHTML = '';
    pagination.classList.add('is-hidden');
    setContentEmpty('taskContentArea', false);
    list.appendChild(el('div', { class: 'empty', html: `${icon('i-clock', 'ic')}<p>Загрузка…</p>` }));

    const pageSize = state.taskPageSize;
    const params = {
        limit: pageSize,
        offset: (state.tasksPage - 1) * pageSize,
    };

    let page;
    try {
        page = await api.listTasks(params);
    } catch (err) {
        list.innerHTML = '';
        toastErr('Не удалось загрузить задачи', err.message);
        setContentEmpty('taskContentArea', true);
        list.appendChild(emptyState('i-inbox', 'Задачи недоступны', 'Проверьте, что сервер API запущен'));
        return;
    }

    const tasks = page.items || [];
    const total = page.total ?? tasks.length;
    const pages = totalPages(total, pageSize);

    if (state.tasksPage > pages) {
        state.tasksPage = pages;
        return renderTasks();
    }

    list.innerHTML = '';
    if (!tasks.length) {
        setContentEmpty('taskContentArea', true);
        list.appendChild(emptyState('i-inbox', 'Задач пока нет', 'Нажмите «Новая задача», чтобы создать первую'));
    } else {
        setContentEmpty('taskContentArea', false);
        tasks.forEach((t) => list.appendChild(taskCard(t)));
    }

    renderPaginationBar(pagination, {
        page: state.tasksPage,
        total,
        pageSize,
        selectId: 'taskPageSize',
        onPage: (p) => {
            state.tasksPage = p;
            renderTasks();
        },
        onPageSize: (size) => {
            state.taskPageSize = size;
            state.tasksPage = 1;
            renderTasks();
        },
    });
}

function taskCard(t) {
    const author = userById(t.author_user_id);
    const authorName = author ? author.full_name : `ID ${t.author_user_id}`;

    const status = t.completed
        ? `<span class="badge-status status-done">${icon('i-check')} Выполнено</span>`
        : `<span class="badge-status status-wip">${icon('i-clock')} В работе · ${relativeAge(t.created_at)}</span>`;

    const descHtml = t.description
        ? `<div class="task-desc">${escapeHtml(t.description)}</div>` : '';

    const card = el('div', { class: `task-card ${t.completed ? 'is-done' : ''}` });
    card.innerHTML = `
        <button class="task-check ${t.completed ? 'is-done' : ''}" title="Переключить статус">${icon('i-check')}</button>
        <div class="task-main">
            <div class="task-title">${escapeHtml(t.title)}</div>
            ${descHtml}
            <div class="chips">
                <span class="chip">${icon('i-user')} <b>${escapeHtml(authorName)}</b></span>
                <span class="chip">${icon('i-calendar')} ${formatDate(t.created_at)}</span>
                ${status}
                <span class="chip chip-plain">v${t.version}</span>
            </div>
        </div>
        <div class="task-actions">
            <button class="icon-btn" data-act="edit" title="Изменить">${icon('i-edit')}</button>
            <button class="icon-btn danger" data-act="del" title="Удалить">${icon('i-trash')}</button>
        </div>`;

    $('.task-check', card).addEventListener('click', () => toggleTask(t));
    $('[data-act="edit"]', card).addEventListener('click', () => openTaskModal(t));
    $('[data-act="del"]', card).addEventListener('click', () => removeTask(t));
    return card;
}

async function toggleTask(t) {
    try {
        const completed = !t.completed;
        await api.patchTask(t.id, { completed });
        if (completed) toastOk('Задача выполнена');
        else toastOk('Задача в работе');
        renderTasks();
    } catch (err) {
        toastErr('Не удалось обновить задачу', err.message);
    }
}

async function removeTask(t) {
    if (!confirm(`Удалить задачу «${t.title}»?`)) return;
    try {
        await api.deleteTask(t.id);
        toastOk('Задача удалена');
        renderTasks();
    } catch (err) {
        toastErr('Не удалось удалить задачу', err.message);
    }
}

/* =========================================================================
 * Раздел: ПОЛЬЗОВАТЕЛИ
 * ========================================================================= */
async function renderUsers() {
    const grid = $('#userGrid');
    const pagination = $('#userPagination');
    grid.innerHTML = '';
    pagination.classList.add('is-hidden');
    setContentEmpty('userContentArea', false);
    grid.appendChild(el('div', { class: 'empty', html: `${icon('i-clock', 'ic')}<p>Загрузка…</p>` }));

    const pageSize = state.userPageSize;
    let page;

    try {
        page = await api.listUsers({
            limit: pageSize,
            offset: (state.usersPage - 1) * pageSize,
        });
    } catch (err) {
        toastErr('Не удалось загрузить пользователей', err.message);
        grid.innerHTML = '';
        setContentEmpty('userContentArea', true);
        grid.appendChild(emptyState('i-users', 'Пользователи недоступны', 'Проверьте, что сервер API запущен'));
        return;
    }

    const users = page.items || [];
    const total = page.total ?? users.length;
    const pages = totalPages(total, pageSize);

    if (state.usersPage > pages) {
        state.usersPage = pages;
        return renderUsers();
    }

    grid.innerHTML = '';
    if (!users.length) {
        setContentEmpty('userContentArea', true);
        grid.appendChild(emptyState('i-users', isAdmin() ? 'Пользователей пока нет' : 'Профиль не найден', isAdmin() ? 'Создайте первого пользователя через регистрацию' : 'Попробуйте сменить пользователя'));
    } else {
        setContentEmpty('userContentArea', false);
        users.forEach((u) => grid.appendChild(userCard(u)));
    }

    if (isAdmin() && total > pageSize) {
        pagination.hidden = false;
        renderPaginationBar(pagination, {
            page: state.usersPage,
            total,
            pageSize,
            selectId: 'userPageSize',
            onPage: (p) => {
                state.usersPage = p;
                renderUsers();
            },
            onPageSize: (size) => {
                state.userPageSize = size;
                state.usersPage = 1;
                renderUsers();
            },
        });
    } else {
        pagination.hidden = true;
        pagination.classList.add('is-hidden');
    }
}

function userCard(u) {
    const phoneHtml = u.phone_number
        ? `<div class="user-phone">${icon('i-phone')} ${escapeHtml(u.phone_number)}</div>`
        : '<div class="user-phone user-phone--placeholder" aria-hidden="true"></div>';

    const card = el('div', { class: 'user-card' });
    card.innerHTML = `
        <div class="user-card-top">
            <div class="user-avatar" style="background:${avatarColor(u.id)}">${escapeHtml(initials(u.full_name))}</div>
            <div class="user-card-body">
                <div class="user-name">${escapeHtml(u.full_name)}</div>
                ${phoneHtml}
            </div>
        </div>
        <div class="user-badges">
            <span class="chip">ID: ${u.id}</span>
            <span class="chip">v${u.version}</span>
        </div>
        <div class="user-actions">
            <button class="btn btn-ghost btn-sm" data-act="edit">${icon('i-edit')} Изменить</button>
            <button class="btn btn-danger-soft btn-sm" data-act="del">${icon('i-trash')} Удалить</button>
        </div>`;

    $('[data-act="edit"]', card).addEventListener('click', () => openUserModal(u));
    $('[data-act="del"]', card).addEventListener('click', () => removeUser(u));
    return card;
}

async function removeUser(u) {
    if (!confirm(`Удалить пользователя «${u.full_name}»?`)) return;
    try {
        await api.deleteUser(u.id);
        toastOk('Пользователь удалён');
        if (u.id === getActingUserId()) {
            clearActingUserId();
            state.actingUser = null;
            updateBrandBadge();
            openIdentityModal();
            return;
        }
        await refreshUsersCatalog();
        renderUsers();
    } catch (err) {
        const hasTasks = err.status === 409 || err.code === 'conflict';
        if (hasTasks) {
            toastErr('Нельзя удалить пользователя', 'Сначала удалите его задачи');
        } else {
            toastErr('Не удалось удалить пользователя', err.message);
        }
    }
}

/* =========================================================================
 * Раздел: СТАТИСТИКА
 * ========================================================================= */
async function renderStats() {
    const params = {
        from: $('#statFrom').value || undefined,
        to: $('#statTo').value || undefined,
    };
    if (isAdmin()) {
        params.user_id = $('#statUserFilter').value || undefined;
    }

    let s;
    try {
        s = await api.statistics(params);
    } catch (err) {
        toastErr('Не удалось загрузить статистику', err.message);
        return;
    }

    const created = s.tasks_created || 0;
    const completed = s.tasks_completed || 0;
    const wip = Math.max(0, created - completed);
    // Бэкенд отдаёт процент выполнения (0–100), а не долю.
    const rate = s.task_completed_rate != null ? s.task_completed_rate : 0;

    $('#statCreated').textContent = created;
    $('#statCompleted').textContent = completed;
    $('#statRate').innerHTML = `${rate.toFixed(1)}<span class="stat-unit">%</span>`;
    $('#statAvg').textContent = formatDuration(s.tasks_average_completion_time);

    const donePct = created ? (completed / created) * 100 : 0;
    const wipPct = created ? (wip / created) * 100 : 0;
    $('#progDoneText').textContent = `${completed} / ${created} (${donePct.toFixed(1)}%)`;
    $('#progWipText').textContent = `${wip} / ${created} (${wipPct.toFixed(1)}%)`;
    $('#progDoneBar').style.width = `${donePct}%`;
    $('#progWipBar').style.width = `${wipPct}%`;

    $('#statUpdated').textContent = `Обновлено: ${new Date().toLocaleTimeString('ru-RU')}`;
}

/* =========================================================================
 * Пустое состояние
 * ========================================================================= */
function emptyState(iconId, title, text) {
    const node = el('div', { class: 'empty' });
    node.innerHTML = `
        <div class="empty-icon">${icon(iconId)}</div>
        <h3>${escapeHtml(title)}</h3>
        <p>${escapeHtml(text || '')}</p>`;
    return node;
}

/* =========================================================================
 * Модалки
 * ========================================================================= */
function openModal(id) {
    closeAllCustomSelects();
    $(`#${id}`).classList.remove('is-hidden');
}
function closeModal(id) { $(`#${id}`).classList.add('is-hidden'); }

const formErrorTimers = new Map();

function hideError(id) {
    const e = $(`#${id}`);
    if (!e) return;
    e.classList.add('is-hidden');
    e.textContent = '';
    const timer = formErrorTimers.get(id);
    if (timer) {
        clearTimeout(timer);
        formErrorTimers.delete(id);
    }
}

function showError(id, msg) {
    const e = $(`#${id}`);
    e.textContent = msg;
    e.classList.remove('is-hidden');
    const prev = formErrorTimers.get(id);
    if (prev) clearTimeout(prev);
    formErrorTimers.set(id, setTimeout(() => hideError(id), 3500));
}

/* --- Задача --- */
async function openTaskModal(task = null) {
    state.editingTaskId = task ? task.id : null;
    hideError('taskFormError');

    $('#taskModalTitle').textContent = task ? 'Редактировать задачу' : 'Новая задача';
    $('#taskSubmit').textContent = task ? 'Сохранить' : 'Создать задачу';
    $('#taskTitle').value = task ? task.title : '';
    $('#taskDescription').value = task && task.description ? task.description : '';

    openModal('taskModal');
    setTimeout(() => $('#taskTitle').focus(), 50);
}

async function submitTask() {
    const title = $('#taskTitle').value.trim();
    const description = $('#taskDescription').value.trim();
    hideError('taskFormError');

    if (!title) return showError('taskFormError', 'Заголовок обязателен.');

    const actingId = getActingUserId();
    if (!actingId) return showError('taskFormError', 'Сначала выберите пользователя.');

    const btn = $('#taskSubmit');
    btn.disabled = true;
    try {
        if (state.editingTaskId == null) {
            const body = { title, author_user_id: actingId };
            if (description) body.description = description;
            await api.createTask(body);
            state.tasksPage = 1;
            toastOk('Задача создана');
        } else {
            const body = { title, description: description || null };
            await api.patchTask(state.editingTaskId, body);
            toastOk('Задача обновлена');
        }
        closeTaskModal();
        renderTasks();
    } catch (err) {
        showError('taskFormError', err.message);
    } finally {
        btn.disabled = false;
    }
}

function closeTaskModal() {
    state.editingTaskId = null;
    closeModal('taskModal');
}

/* --- Пользователь --- */
function openUserModal(user = null) {
    const profile = user || state.actingUser;
    if (!profile) return;

    state.editingUserId = profile.id;
    hideError('userFormError');

    $('#userModalTitle').textContent = 'Редактировать профиль';
    $('#userSubmit').textContent = 'Сохранить';
    $('#userFullName').value = profile.full_name;
    $('#userPhone').value = profile.phone_number ? profile.phone_number : '';
    $('#userClearPhone').checked = false;
    $('#userClearPhoneField').classList.toggle('is-hidden', false);

    openModal('userModal');
    setTimeout(() => $('#userFullName').focus(), 50);
}

async function submitUser() {
    const fullName = $('#userFullName').value.trim();
    const phone = $('#userPhone').value.trim();
    const clearPhone = $('#userClearPhone').checked;
    hideError('userFormError');

    if (fullName.length < 3) return showError('userFormError', 'Полное имя должно содержать минимум 3 символа.');
    if (phone && !phone.startsWith('+')) return showError('userFormError', 'Телефон должен начинаться с «+».');
    if (phone && (phone.length < 10 || phone.length > 15)) {
        return showError('userFormError', 'Телефон должен быть от 10 до 15 символов.');
    }

    const btn = $('#userSubmit');
    btn.disabled = true;
    try {
        const body = { full_name: fullName };
        if (clearPhone) body.phone_number = null;
        else if (phone) body.phone_number = phone;
        await api.patchUser(state.editingUserId, body);
        toastOk('Профиль обновлён');
        closeModal('userModal');
        await refreshUsersCatalog();
        renderUsers();
    } catch (err) {
        showError('userFormError', err.message);
    } finally {
        btn.disabled = false;
    }
}

/* =========================================================================
 * Тема
 * ========================================================================= */
const THEME_KEY = 'todoapp_theme';
function applyTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    const dark = theme === 'dark';
    $('#themeLabel').textContent = dark ? 'Light mode' : 'Dark mode';
    $('#themeIcon').innerHTML = `<use href="#${dark ? 'i-sun' : 'i-moon'}"/>`;
    localStorage.setItem(THEME_KEY, theme);
}
function toggleTheme() {
    const current = document.documentElement.getAttribute('data-theme');
    const next = current === 'dark' ? 'light' : 'dark';
    applyTheme(next);
    toastOk(next === 'dark' ? 'Тема переключена на тёмную' : 'Тема переключена на светлую');
}

/* =========================================================================
 * Ссылка на Swagger
 * ========================================================================= */
function updateApiLink() {
    try { $('#apiHost').textContent = new URL(API_BASE).host; }
    catch { $('#apiHost').textContent = API_BASE; }
    const link = $('#apiSwaggerLink');
    if (link) link.href = `${API_BASE}/swagger/`;
}

/* =========================================================================
 * Приветственный занавес
 * ========================================================================= */
const WELCOME_HOLD_MS = 3200;
const WELCOME_OPEN_MS = 1800;

function playWelcomeCurtain() {
    const curtain = $('#welcomeCurtain');
    if (!curtain) return;

    document.body.classList.add('welcome-active');

    window.requestAnimationFrame(() => {
        window.requestAnimationFrame(() => curtain.classList.add('is-ready'));
    });

    const finish = () => {
        curtain.remove();
        document.body.classList.remove('welcome-active');
    };

    window.setTimeout(() => {
        curtain.classList.add('is-opening');
        const panel = curtain.querySelector('.welcome-curtain__panel--left');
        if (panel) {
            panel.addEventListener('transitionend', finish, { once: true });
        }
        window.setTimeout(finish, WELCOME_OPEN_MS + 120);
    }, WELCOME_HOLD_MS);
}

function startApp() {
    applyTheme(localStorage.getItem(THEME_KEY) || 'dark');
    updateApiLink();
    initCustomSelects();
    bindEvents();
    updateRoleUI();

    if (!getActingUserId() && !isAdmin()) {
        openIdentityModal();
        return;
    }

    loadActingUser()
        .catch((err) => {
            clearActingSession();
            toastErr('Сессия недействительна', err.message);
            openIdentityModal();
        })
        .finally(() => switchView(isAdmin() ? 'users' : 'tasks'));
}

/* =========================================================================
 * Инициализация
 * ========================================================================= */
function isModalOpen(id) {
    const backdrop = $(`#${id}`);
    return backdrop && !backdrop.classList.contains('is-hidden');
}

function isCustomSelectOpen() {
    return $$('.cselect.is-open').length > 0;
}

function handleModalEnter(e) {
    if (e.key !== 'Enter' || e.isComposing) return;
    if (isCustomSelectOpen()) return;

    const tag = e.target.tagName;
    if (tag === 'BUTTON') return;

    if (tag === 'TEXTAREA') {
        if (e.shiftKey) return; // Shift+Enter — новая строка
    }

    if (isModalOpen('taskModal')) {
        e.preventDefault();
        submitTask();
        return;
    }

    if (isModalOpen('userModal')) {
        e.preventDefault();
        submitUser();
        return;
    }

    if (isModalOpen('identityModal') && state.identityAdminMode) {
        e.preventDefault();
        submitIdentityAdmin();
        return;
    }

    if (isModalOpen('identityModal') && state.identityRegisterMode) {
        e.preventDefault();
        submitIdentityRegister();
    }
}

function bindEvents() {
    $$('.nav-item').forEach((b) => b.addEventListener('click', () => switchView(b.dataset.view)));

    $('#themeToggle').addEventListener('click', toggleTheme);
    $('#switchUserBtn').addEventListener('click', switchActingUser);
    $('#identityToggleRegister').addEventListener('click', () => setIdentityRegisterMode(true));
    $('#identityToggleAdmin').addEventListener('click', () => setIdentityAdminMode(true));
    $('#identityBackBtn').addEventListener('click', () => backToIdentityList({ refreshList: true }));
    $('#identityCancelBtn').addEventListener('click', closeIdentityModal);
    $('#identityCloseBtn').addEventListener('click', handleIdentityDismiss);
    $('#identitySubmitRegister').addEventListener('click', submitIdentityRegister);
    $('#identitySubmitAdmin').addEventListener('click', submitIdentityAdmin);

    // Задачи
    $('#newTaskBtn').addEventListener('click', () => openTaskModal(null));
    $('#taskSubmit').addEventListener('click', submitTask);

    // Профиль
    $('#userSubmit').addEventListener('click', submitUser);

    // Статистика
    $('#statApply').addEventListener('click', renderStats);
    $('#statReset').addEventListener('click', () => {
        if (isAdmin()) setSelectValue($('#statUserFilter'), '');
        $('#statFrom').value = '';
        $('#statTo').value = '';
        renderStats();
    });

    // Закрытие модалок — только по кнопке «Отмена»/крестику или Escape (не по клику на фон)
    $$('[data-close-modal]').forEach((b) =>
        b.addEventListener('click', () => {
            if (b.dataset.closeModal === 'taskModal') closeTaskModal();
            else closeModal(b.dataset.closeModal);
        }));
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            if (isModalOpen('identityModal')) {
                handleIdentityDismiss();
                return;
            }
            $$('.modal-backdrop').forEach((bd) => {
                if (bd.classList.contains('is-hidden')) return;
                if (bd.id === 'taskModal') closeTaskModal();
                else bd.classList.add('is-hidden');
            });
            return;
        }

        handleModalEnter(e);
    });
}

async function init() {
    startApp();
}

playWelcomeCurtain();

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
} else {
    init();
}
