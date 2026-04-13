import * as bindings from "../../bindings/changeme/internal/app/context.js";

type PopupBindings = {
  OpenAbout: () => Promise<void>;
};

const aboutDialogBindings = bindings as unknown as PopupBindings;

export function openAboutDialog() {
  return aboutDialogBindings.OpenAbout();
}
