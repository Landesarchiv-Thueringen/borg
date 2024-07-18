import { CommonModule } from '@angular/common';
import { Component, Inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatTableModule } from '@angular/material/table';
import { FileFeaturePipe } from '../pipes/file-attribut-de.pipe';
import { PrettyPrintCsvPipe } from '../pipes/pretty-print-csv.pipe';
import { PrettyPrintJsonPipe } from '../pipes/pretty-print-json.pipe';
import { ToolResult } from '../results';

interface DialogData {
  toolResult: ToolResult;
}

@Component({
  selector: 'app-tool-output',
  templateUrl: './tool-output.component.html',
  styleUrls: ['./tool-output.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    FileFeaturePipe,
    MatButtonModule,
    MatDialogModule,
    MatExpansionModule,
    MatTableModule,
    PrettyPrintCsvPipe,
    PrettyPrintJsonPipe,
  ],
})
export class ToolOutputComponent {
  readonly toolResult = this.data.toolResult;

  constructor(@Inject(MAT_DIALOG_DATA) private data: DialogData) {}
}
