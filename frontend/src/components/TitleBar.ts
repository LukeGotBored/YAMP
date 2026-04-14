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
    label?: string;
    ariaLabel?: string;
    view: any;
    placement?: "left" | "right";
    icon?: Component<any>;
    iconSize?: number;
    iconWeight?: "thin" | "light" | "regular" | "bold" | "fill" | "duotone";
    activeIconWeight?: "thin" | "light" | "regular" | "bold" | "fill" | "duotone";
    hidden?: boolean;
}

export function tabButtonId(tabId: string) {
    return `main-tab-${tabId}`;
}

export function tabPanelId(tabId: string) {
    return `main-panel-${tabId}`;
}
