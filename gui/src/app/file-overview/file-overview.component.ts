// angular
import { Component, Inject } from '@angular/core';

// material
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';

// project
import { FileAnalysisService, Summary, ToolConfidence } from '../file-analysis/file-analysis.service';
import { FeatureValue, FileResult, ToolResult } from '../file-analysis/file-analysis.service';

interface DialogData {
  fileResult: FileResult;
}

interface FileFeature {
  value: string;
  confidence?: number;
  colorizeConfidence?: boolean;
}

interface FileFeatures {
  [key: string]: FileFeature;
}

@Component({
  selector: 'app-file-overview',
  templateUrl: './file-overview.component.html',
  styleUrls: ['./file-overview.component.scss'],
})
export class FileOverviewComponent {
  readonly fileResult: FileResult;
  dataSource: MatTableDataSource<FileFeatures>;
  tableColumnList: string[];
  constructor(
    @Inject(MAT_DIALOG_DATA) private data: DialogData,
    private fileAnalysisService: FileAnalysisService,
  ) {
    this.dataSource = new MatTableDataSource<FileFeatures>();
    this.tableColumnList = ['Attribut'];
    this.fileResult = data.fileResult;
    this.initTableData();
  }

  initTableData(): void {
    if (this.fileResult.toolResults.summary) {
      const summary = this.fileResult.toolResults.summary;
      const toolNames: string[] = [];
      let featureNames: string[] = [];
      let toolResults: ToolResult[] = this.fileResult.toolResults.fileIdentificationResults;
      if (this.fileResult.toolResults.fileValidationResults) {
        toolResults = toolResults.concat(this.fileResult.toolResults.fileValidationResults);
      }
      toolResults.forEach(
        (toolResult: ToolResult) => {
          toolNames.push(toolResult.toolName);
        }
      );
      for (let featureKey in summary) {
        featureNames.push(featureKey);
      }
      featureNames = this.fileAnalysisService.selectOverviewFeatures(featureNames);
      this.dataSource.data = this.getTableRows(summary, toolNames, featureNames);
    }
  }

  getTableRows(summary: Summary, toolNames: string[], featureNames: string[]): FileFeatures[] {
    const rows: FileFeatures[] = [this.getCumulativeResult(summary, featureNames)];
    const sortedFeatures: string[] = this.fileAnalysisService.sortFeatures(featureNames);
    this.tableColumnList = ['Werkzeug', ...sortedFeatures];
    for (let toolName of toolNames) {
      const featureValues: FileFeatures = {};
      featureValues['Werkzeug'] = {
        value: toolName,
      }
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
      rows.push(featureValues);
    }
    return rows;
  }

  getCumulativeResult(summary: Summary, featureNames: string[]): FileFeatures {
    const features: FileFeatures = {};
    features['Werkzeug'] = {
      value: 'Gesamtergebnis',
    }
    for (let featureName of featureNames) {
      // result with highest confidence
      const featureValue = summary[featureName].values[0]
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
}
