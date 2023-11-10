// angular
import { AfterViewInit, Component, ViewChild } from '@angular/core';

// material
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';

// project
import {
  FileResult,
  FileAnalysisService,
  Feature,
} from '../file-analysis.service';
import { FileSizePipe } from '../../utility/file-size/file-size.pipe';
import { FileOverviewComponent } from 'src/app/file-overview/file-overview.component';

export interface FileOverview {
  [key: string]: FileFeature;
}

export interface FileFeature {
  value: string;
  feature?: Feature;
  tooltip?: string;
}

@Component({
  selector: 'app-file-analysis-table',
  templateUrl: './file-analysis-table.component.html',
  styleUrls: ['./file-analysis-table.component.scss'],
})
export class FileAnalysisTableComponent implements AfterViewInit {
  dataSource: MatTableDataSource<FileOverview>;
  generatedTableColumnList: string[];
  tableColumnList: string[];

  @ViewChild(MatPaginator) paginator!: MatPaginator;

  constructor(
    private dialog: MatDialog,
    private fileAnalysisService: FileAnalysisService,
    private fileSizePipe: FileSizePipe,
  ) {
    this.dataSource = new MatTableDataSource<FileOverview>([]);
    this.tableColumnList = [];
    this.generatedTableColumnList = ['fileName', 'relativePath', 'fileSize'];
    this.fileAnalysisService.getFileResults().subscribe({
      // error can't occure --> no error handling
      next: (fileInfos: FileResult[]) => {
        this.processFileInformations(fileInfos);
      },
    });
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
  }

  processFileInformations(fileInfos: FileResult[]): void {
    const featureKeys: string[] = ['fileName', 'relativePath', 'fileSize'];
    const data: FileOverview[] = [];
    for (let fileInfo of fileInfos) {
      let fileOverview: FileOverview = {};
      for (let featureKey in fileInfo.toolResults.summary) {
        featureKeys.push(featureKey);
        fileOverview['fileName'] = { value: fileInfo.fileName };
        fileOverview['relativePath'] = fileInfo.relativePath
          ? { value: fileInfo.relativePath }
          : { value: '' };
        fileOverview['fileSize'] = {
          value: this.fileSizePipe.transform(fileInfo.fileSize),
        };
        fileOverview[featureKey] = {
          value: fileInfo.toolResults.summary[featureKey].values[0].value,
          feature: fileInfo.toolResults.summary[featureKey],
          tooltip: this.getFeatureTooltip(
            fileInfo.toolResults.summary[featureKey]
          ),
        };
      }
      data.push(fileOverview);
    }
    this.dataSource.data = data;
    const features = [...new Set(featureKeys)];
    const selectedFeatures = this.selectFeatures(features);
    const sortedFeatures = this.sortFeatures(selectedFeatures);
    this.generatedTableColumnList = sortedFeatures;
    this.tableColumnList = sortedFeatures.concat(['actions']);
  }

  selectFeatures(features: string[]): string[] {
    const overviewFeatures: string[] = this.fileAnalysisService.getOverviewFeatures();
    return features.filter((feature: string) => {
      return overviewFeatures.includes(feature);
    });
  }

  sortFeatures(features: string[]): string[] {
    return features.sort((f1: string, f2: string) => {
      const featureOrder = this.fileAnalysisService.getFeatureOrder();
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

  getFeatureTooltip(feature: Feature): string {
    let tooltip = '';
    for (let featureValue of feature.values) {
      tooltip +=
        '[' + featureValue.value + '; ' + featureValue.score.toFixed(2) + ']: ';
      for (let tool of featureValue.tools) {
        tooltip +=
          '(' + tool.toolName + '; ' + tool.confidence.toFixed(2) + ')';
      }
      tooltip += '\n';
    }
    return tooltip;
  }

  openDetails(fileOverview: FileOverview): void {
    this.dialog.open(FileOverviewComponent, {
      data: {
        fileOverview: fileOverview
      }
    });
  }
}
