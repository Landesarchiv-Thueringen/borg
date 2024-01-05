import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToolResult } from '../file-analysis/file-analysis.service';

interface DialogData {
  toolResult: ToolResult;
}

@Component({
  selector: 'app-tool-output',
  templateUrl: './tool-output.component.html',
  styleUrls: ['./tool-output.component.scss'],
})
export class ToolOutputComponent {
  readonly toolResult = this.data.toolResult;

  constructor(@Inject(MAT_DIALOG_DATA) private data: DialogData) {}
}
