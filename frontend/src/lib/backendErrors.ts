export const BACKEND_FATAL_ERROR_EVENT = "app:fatal-error";
export const BACKEND_PANIC_EVENT = "app:panic";

export function formatBackendError(prefix: string, eventData: any): string {
    const payload = Array.isArray(eventData) ? eventData[0] : eventData;
    
    if (!payload) return `${prefix} Unknown backend error occurred.\nPayload: ${JSON.stringify(eventData)}`;
    
    // Sometimes the Wails event wraps the payload in an object containing the emitted data
    const safeMsg = payload.Message || payload.message || payload.error || "Unknown errored payload";
    const safeTime = payload.Time || payload.time || new Date().toISOString();
    const safeTrace = payload.StackTrace || payload.stackTrace || "No trace available (or trace was skipped)";
    
    return `=========================================
${prefix}
=========================================
Time:      ${safeTime}
Message:   ${safeMsg}

--- Environment ---
UserAgent: ${navigator.userAgent}
Location:  ${window.location.href}
Viewport:  ${window.innerWidth}x${window.innerHeight}

--- Backend Stack Trace ---
${safeTrace}
=========================================`;
}
