// angular
import { Component } from '@angular/core';

// material
import { MatTableDataSource } from '@angular/material/table';

// project
import { FileInformation, FileAnalysisService } from '../file-analysis.service';

export interface FileOverview {
  [key: string]: string
}

@Component({
  selector: 'app-file-analysis-table',
  templateUrl: './file-analysis-table.component.html',
  styleUrls: ['./file-analysis-table.component.scss']
})
export class FileAnalysisTableComponent {
  dataSource: MatTableDataSource<FileOverview>;
  tableColumnList: string[];

  constructor(private fileAnalysisService: FileAnalysisService) {
    this.dataSource = new MatTableDataSource<FileOverview>([]);
    this.tableColumnList = ['fileName', 'relativePath', 'fileSize'];
    this.fileAnalysisService.getFileInfo().subscribe({
      // error can't occure --> no error handling
      next: (fileInfos: FileInformation[]) => {
        console.log(fileInfos);
        this.processFileInformations(fileInfos);
      }
    })
  }

  processFileInformations(fileInfos: FileInformation[]): void {
    const featureKeys: string[] = ['fileName', 'relativePath', 'fileSize']
    const data: FileOverview[] = []
    for (let fileInfo of fileInfos) {
      let fileOverview: FileOverview = {}
      for (let featureKey in fileInfo.fileAnalysis.summary) {
        featureKeys.push(featureKey)
        fileOverview['fileName'] = fileInfo.fileName
        fileOverview['relativePath'] = fileInfo.relativePath ? fileInfo.relativePath : ''
        fileOverview['fileSize'] = fileInfo.size
        fileOverview[featureKey] = fileInfo.fileAnalysis.summary[featureKey].values[0].value
      }
      data.push(fileOverview)
    }
    this.dataSource.data = data
    this.tableColumnList = [...new Set(featureKeys)]
  }
}
