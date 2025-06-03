import { CommonModule } from '@angular/common';
import { Component, inject, OnInit } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatRippleModule } from '@angular/material/core';
import { MAT_DIALOG_DATA, MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatTabsModule } from '@angular/material/tabs';
import { RouterModule } from '@angular/router';
import { FeatureValuePipe } from '../pipes/feature-value.pipe';
import { FeatureSet, FileAnalysis, ToolFeatureValue, ToolResult } from '../results';
import { ToolOutputComponent } from '../tool-output/tool-output.component';

export interface DialogData {
  featureSet: FeatureSet;
  toolResults: ToolResult[];
  analysis: FileAnalysis;
}

interface ToolRow {
  toolName: string;
  puid: ToolFeatureValue | undefined;
  mimeType: ToolFeatureValue | undefined;
  formatVersion: ToolFeatureValue | undefined;
  valid: ToolFeatureValue | undefined;
  error: boolean;
}

@Component({
  selector: 'app-result-details',
  imports: [
    CommonModule,
    MatButtonModule,
    MatDialogModule,
    MatIconModule,
    MatTableModule,
    RouterModule,
    MatTabsModule,
    FeatureValuePipe,
    MatRippleModule,
  ],
  templateUrl: './result-details.component.html',
  styleUrl: './result-details.component.scss',
})
export class ResultDetailsComponent implements OnInit {
  private readonly data = inject<DialogData>(MAT_DIALOG_DATA);
  private readonly dialog = inject(MatDialog);
  private readonly featureSet: FeatureSet = this.data.featureSet;
  private readonly toolResults: ToolResult[] = this.data.toolResults;
  readonly rows: ToolRow[] = [];
  displayedColumns: string[] = ['tool', 'puid', 'mimeType', 'formatVersion', 'valid'];

  ngOnInit() {
    this.rows.push({
      toolName: 'Gesamtergebnis',
      puid: this.featureSet.features['format:puid'],
      mimeType: this.featureSet.features['format:mimeType'],
      formatVersion: this.featureSet.features['format:version'],
      valid: this.featureSet.features['format:valid'],
      error: false,
    });
    for (let toolResult of this.toolResults) {
      if (this.featureSet.supportingTools.includes(toolResult.id) || toolResult.error) {
        this.rows.push({
          toolName: toolResult.title,
          puid: toolResult.features['format:puid'],
          mimeType: toolResult.features['format:mimeType'],
          formatVersion: toolResult.features['format:version'],
          valid: toolResult.features['format:valid'],
          error: !!toolResult.error,
        });
      }
    }
    if (this.rows.some((row) => row.error)) {
      this.displayedColumns.push('error');
    }
  }

  showToolOutput(toolName: string): void {
    const toolResult = this.toolResults.find((r) => r.title === toolName);
    if (toolResult) {
      this.dialog.open(ToolOutputComponent, {
        data: {
          toolName,
          toolResult,
        },
        autoFocus: false,
        maxWidth: '80vw',
      });
    }
  }
}
