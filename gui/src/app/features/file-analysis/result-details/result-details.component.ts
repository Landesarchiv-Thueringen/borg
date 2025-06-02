import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { MatTabsModule } from '@angular/material/tabs';
import { RouterModule } from '@angular/router';
import { FeatureSetsTableComponent } from '../feature-sets-table/feature-sets-table.component';
import { FileFeaturePipe } from '../pipes/file-feature.pipe';
import { FileAnalysis, RowValue } from '../results';
import { ToolOutputComponent } from '../tool-output/tool-output.component';

const OVERVIEW_FEATURES = [
  'format:puid',
  'format:mimeType',
  'format:version',
  'format:valid',
] as const;

const featureOrder = new Map<string, number>([
  ['puid', 4],
  ['mimeType', 5],
  ['formatVersion', 6],
  ['encoding', 7],
  ['', 101],
  ['wellFormed', 1001],
  ['valid', 1002],
]);

interface DialogData {
  filename: string;
  info: { [key: string]: RowValue };
  analysis: FileAnalysis;
}

interface FileFeature {
  value?: string | boolean | number;
  confidence?: number;
  icon?: string;
}

interface FileFeatures {
  [key: string]: FileFeature;
}

@Component({
  selector: 'app-result-details',
  imports: [
    CommonModule,
    FileFeaturePipe,
    MatButtonModule,
    MatDialogModule,
    MatIconModule,
    MatTableModule,
    RouterModule,
    MatTabsModule,
  ],
  templateUrl: './result-details.component.html',
  styleUrl: './result-details.component.scss',
})
export class ResultDetailsComponent {
  data = inject<DialogData>(MAT_DIALOG_DATA);
  private dialog = inject(MatDialog);
  readonly analysis: FileAnalysis = this.data.analysis;
  dataSource = new MatTableDataSource<FileFeatures>();
  tableColumnList: string[] = [];

  constructor() {
    this.initTableData();
  }

  initTableData(): void {
    if (this.analysis.featureSets.length === 0) {
      return;
    }
    const sortedFeatures: string[] = sortFeatures([...OVERVIEW_FEATURES]);
    this.tableColumnList = ['tool', ...sortedFeatures];
    if (this.analysis.summary.error) {
      this.tableColumnList.push('error');
    }
    const rows: FileFeatures[] = [];
    const row: FileFeatures = {};
    row['tool'] = {
      value: 'Gesamtergebnis',
    };
    for (let featureName of OVERVIEW_FEATURES) {
      if (this.analysis.featureSets[0].features[featureName] !== undefined) {
        row[featureName] = {
          value: this.analysis.featureSets[0].features[featureName].value,
        };
      }
    }
    rows.push(row);
    for (let toolResult of this.analysis.toolResults) {
      const row: FileFeatures = {};
      row['tool'] = { value: toolResult.title };
      for (let featureName of OVERVIEW_FEATURES) {
        if (toolResult.features[featureName] !== undefined) {
          row[featureName] = {
            value: toolResult.features[featureName].value,
          };
        }
      }
      if (toolResult.error) {
        row['error'] = {
          icon: 'error',
        };
      }
      rows.push(row);
    }
    this.dataSource.data = rows;
  }

  showToolOutput(toolName: string): void {
    const toolResult = this.analysis.toolResults.find((r) => r.title === toolName);
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

  showFeatureSets(): void {
    this.dialog.open(FeatureSetsTableComponent, {
      data: {
        featureSets: this.analysis.featureSets,
        toolResults: this.analysis.toolResults,
      },
      autoFocus: false,
      maxWidth: '80vw',
    });
  }

  exportResult(): void {
    const a = document.createElement('a');
    document.body.appendChild(a);
    a.download = 'borg-results.json';
    a.href =
      'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(this.analysis, null, 2));
    a.click();
    document.body.removeChild(a);
  }
}

/** Sorts feature keys and removes duplicates. */
function sortFeatures(features: string[]): string[] {
  features = [...new Set(features)];
  return features.sort((f1: string, f2: string) => {
    let orderF1: number | undefined = featureOrder.get(f1);
    if (!orderF1) {
      orderF1 = featureOrder.get('');
    }
    let orderF2: number | undefined = featureOrder.get(f2);
    if (!orderF2) {
      orderF2 = featureOrder.get('');
    }
    if (orderF1! < orderF2!) {
      return -1;
    } else if (orderF1! > orderF2!) {
      return 1;
    }
    return 0;
  });
}
