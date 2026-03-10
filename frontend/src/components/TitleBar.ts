export class ItemTab {
  id: string;
  label: string;
  view: any;
  hidden?: boolean;

  constructor(id: string, label: string, view: any) {
    this.id = id;
    this.label = label;
    this.view = view;
    this.hidden = false;
  }
}
