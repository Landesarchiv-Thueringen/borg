// angular
import { AfterViewInit, Component, ViewChild } from '@angular/core';

// material
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';

// project
import { FileResult, FileAnalysisService } from '../file-analysis.service';
import { FileSizePipe } from '../../utility/file-size/file-size.pipe';

export interface FileOverview {
  [key: string]: string;
}

@Component({
  selector: 'app-file-analysis-table',
  templateUrl: './file-analysis-table.component.html',
  styleUrls: ['./file-analysis-table.component.scss'],
})
export class FileAnalysisTableComponent implements AfterViewInit {
  dataSource: MatTableDataSource<FileOverview>;
  tableColumnList: string[];

  @ViewChild(MatPaginator) paginator!: MatPaginator;

  constructor(
    private fileAnalysisService: FileAnalysisService,
    private fileSizePipe: FileSizePipe,
  ) {
    this.dataSource = new MatTableDataSource<FileOverview>([]);
    this.tableColumnList = ['fileName', 'relativePath', 'fileSize'];
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
        fileOverview['fileName'] = fileInfo.fileName;
        fileOverview['relativePath'] = fileInfo.relativePath
          ? fileInfo.relativePath
          : '';
        fileOverview['fileSize'] = this.fileSizePipe.transform(fileInfo.fileSize);
        fileOverview[featureKey] =
          fileInfo.toolResults.summary[featureKey].values[0].value;
      }
      data.push(fileOverview);
    }
    this.dataSource.data = data;
    const features = [...new Set(featureKeys)];
    this.tableColumnList = this.sortFeatures(features);
  }

  sortFeatures(features: string[]): string[] {
    return features.sort((f1: string, f2: string) => {
      const featureOrder = this.fileAnalysisService.getFeatureOrder();
      let orderF1: number|undefined = featureOrder.get(f1);
      if (!orderF1) {
        orderF1 = featureOrder.get('');
      }
      let orderF2: number|undefined = featureOrder.get(f2);
      if (!orderF2) {
        orderF2 = featureOrder.get('');
      }
      if (orderF1! < orderF2!) {
        return -1;
      } else if (orderF1! > orderF2!) {
        return 1;
      }
      return 0
    });
  }
}
