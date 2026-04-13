import type { Component } from "svelte";

export interface DropdownItem {
    type: "item";
    label: string;
    action: () => void;
    icon?: Component<any>;
    shortcut?: string | string[];
    disabled?: boolean;
}

export interface DropdownSeparator {
    type: "separator";
}

export type DropdownOption = DropdownItem | DropdownSeparator;

export interface ItemTab {
    id: string;
    label: string;
    view: any;
    hidden?: boolean;
}

export function tabButtonId(tabId: string) {
    return `main-tab-${tabId}`;
}

export function tabPanelId(tabId: string) {
    return `main-panel-${tabId}`;
}
