// angular
import { Component, Inject } from '@angular/core';

// material
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';

// project
import { FileAnalysisService } from '../file-analysis/file-analysis.service';
import { FeatureValue, FileResult, ToolResult } from '../file-analysis/file-analysis.service';

export interface DialogData {
  fileResult: FileResult;
}

export interface FileFeatures {
  [key: string]: string;
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
      const rows: FileFeatures[] = [];
      const toolNames: string[] = [];
      const featureNames: string[] = [];
      const toolResults: ToolResult[] =
        this.fileResult.toolResults.fileIdentificationResults.concat(
          this.fileResult.toolResults.fileValidationResults
        );
      toolResults.forEach(
        (toolResult: ToolResult) => {
          toolNames.push(toolResult.toolName);
        }
      );
      for (let featureKey in summary) {
        featureNames.push(featureKey);
      }
      const sortedFeatures: string[] = this.fileAnalysisService.sortFeatures(featureNames);
      this.tableColumnList = ['Werkzeug', ...sortedFeatures];
      for (let toolName of toolNames) {
        const featureValues: FileFeatures = {};
        featureValues['Werkzeug'] = toolName;
        for (let featureName of featureNames) {
          for (let featureValue of summary[featureName].values) {
            if (this.featureOfTool(featureValue, toolName)) {
              let value: string = featureValue.value;
              //const score: string = (featureValue.score * 100).toFixed(2)
              featureValues[featureName] = featureValue.value;
            }
          }
        }
        rows.push(featureValues);
      }
      this.dataSource.data = rows;
      console.log(this.tableColumnList);
    }
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
