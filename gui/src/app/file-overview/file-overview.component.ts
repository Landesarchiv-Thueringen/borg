/* BorgFormat - File format identification and validation
 * Copyright (C) 2024 Landesarchiv Thüringen
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import {
  FeatureValue,
  FileAnalysisService,
  FileResult,
  OverviewFeature,
  Summary,
  ToolConfidence,
  ToolResult,
} from '../file-analysis/file-analysis.service';
import { StatusIconsService } from '../file-analysis/status-icons.service';
import { ToolOutputComponent } from '../tool-output/tool-output.component';

interface DialogData {
  fileResult: FileResult;
}

interface FileFeature {
  value?: string;
  confidence?: number;
  colorizeConfidence?: boolean;
  icon?: string;
}

interface FileFeatures {
  [key: string]: FileFeature;
}

const ALWAYS_VISIBLE_COLUMNS = ['puid', 'mimeType'];

@Component({
  selector: 'app-file-overview',
  templateUrl: './file-overview.component.html',
  styleUrls: ['./file-overview.component.scss'],
})
export class FileOverviewComponent {
  readonly fileResult: FileResult = this.data.fileResult;
  readonly icons = this.statusIcons.getIcons(this.data.fileResult);
  dataSource = new MatTableDataSource<FileFeatures>();
  tableColumnList: string[] = [];

  constructor(
    @Inject(MAT_DIALOG_DATA) private data: DialogData,
    private fileAnalysisService: FileAnalysisService,
    private statusIcons: StatusIconsService,
    private dialog: MatDialog,
  ) {
    this.initTableData();
  }

  initTableData(): void {
    if (this.fileResult.toolResults.summary) {
      const summary = this.fileResult.toolResults.summary;
      const toolNames: string[] = [];
      let featureNames: OverviewFeature[] = [];
      let toolResults: ToolResult[] = this.fileResult.toolResults.fileIdentificationResults ?? [];
      if (this.fileResult.toolResults.fileValidationResults) {
        toolResults = toolResults.concat(this.fileResult.toolResults.fileValidationResults);
      }
      toolResults.forEach((toolResult: ToolResult) => {
        toolNames.push(toolResult.toolName);
      });
      for (let featureKey in summary) {
        if (this.fileAnalysisService.isOverviewFeature(featureKey)) {
          featureNames.push(featureKey);
        }
      }
      this.dataSource.data = this.getTableRows(summary, toolNames, featureNames);
      console.log('table rows', this.dataSource.data);
    }
  }

  getTableRows(summary: Summary, toolNames: string[], featureNames: string[]): FileFeatures[] {
    const rows: FileFeatures[] = [this.getCumulativeResult(summary, featureNames)];
    const sortedFeatures: string[] = this.fileAnalysisService.sortFeatures([
      ...ALWAYS_VISIBLE_COLUMNS,
      ...featureNames,
    ]);
    this.tableColumnList = ['tool', ...sortedFeatures];
    if (this.icons.error) {
      this.tableColumnList.push('error');
    }
    for (let toolName of toolNames) {
      const featureValues: FileFeatures = {};
      featureValues['tool'] = {
        value: toolName,
      };
      for (let featureName of featureNames) {
        for (let featureValue of summary[featureName].values) {
          if (this.featureOfTool(featureValue, toolName)) {
            const toolInfo: ToolConfidence = featureValue.tools.find((toolInfo: ToolConfidence) => {
              return toolInfo.toolName === toolName;
            })!;
            featureValues[featureName] = {
              value: featureValue.value,
              confidence: toolInfo.confidence,
              colorizeConfidence: false,
            };
          }
        }
      }
      if (this.findToolResult(toolName)?.error) {
        featureValues['error'] = {
          icon: 'error',
        };
      }
      rows.push(featureValues);
    }
    return rows;
  }

  getCumulativeResult(summary: Summary, featureNames: string[]): FileFeatures {
    const features: FileFeatures = {};
    features['tool'] = {
      value: 'Gesamtergebnis',
    };
    for (let featureName of featureNames) {
      // result with highest confidence
      const featureValue = summary[featureName].values[0];
      features[featureName] = {
        value: featureValue.value,
        confidence: featureValue.score,
        colorizeConfidence: true,
      };
    }
    return features;
  }

  featureOfTool(featureValue: FeatureValue, toolName: string): boolean {
    for (let tool of featureValue.tools) {
      if (tool.toolName === toolName) {
        return true;
      }
    }
    return false;
  }

  showToolOutput(toolName: string): void {
    const toolResult = this.findToolResult(toolName);
    if (toolResult) {
      this.dialog.open(ToolOutputComponent, {
        data: {
          toolResult,
        },
        autoFocus: false,
      });
    }
  }

  exportResult(): void {
    const a = document.createElement('a');
    document.body.appendChild(a);
    a.download = 'borg-results.json';
    a.href = 'data:text/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(this.fileResult, null, 2));
    a.click();
    document.body.removeChild(a);
  }

  private findToolResult(toolName: string): ToolResult | undefined {
    return [
      ...(this.fileResult.toolResults.fileIdentificationResults ?? []),
      ...(this.fileResult.toolResults.fileValidationResults ?? []),
    ].find((result) => result.toolName === toolName);
  }
}
