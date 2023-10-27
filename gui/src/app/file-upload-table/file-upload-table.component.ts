// angular
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { HttpEventType, HttpEvent } from '@angular/common/http';

// material
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';

// project
import { 
  FileAnalysis, 
  FileInformation, 
  FileAnalysisService,
} from '../file-analysis/file-analysis.service';

export interface FileUpload {
  fileName: string;
  relativePath?: string;
  fileSize: number;
  uploadProgress?: number;
}

@Component({
  selector: 'app-file-upload-table',
  templateUrl: './file-upload-table.component.html',
  styleUrls: ['./file-upload-table.component.scss'],
})
export class FileUploadTableComponent implements AfterViewInit {
  dataSource: MatTableDataSource<FileUpload>;
  displayedColumns: string[];

  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  constructor(private fileAnalysisService: FileAnalysisService) {
    this.dataSource = new MatTableDataSource<FileUpload>();
    this.displayedColumns = [
      'fileName',
      'relativePath',
      'fileSize',
      'uploadProgress',
    ];
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;
  }

  addFile(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const files: FileList | null = input.files;
    if (files && files.length === 1) {
      const file = files[0];
      console.log(file);
      const fileUpload: FileUpload = {
        fileName: file.name,
        fileSize: file.size,
      }
      const data = this.dataSource.data
      const fileDataIndex = data.length;
      data.push(fileUpload)
      this.uploadFile(file, fileDataIndex);
      this.dataSource.data = data
    }
  }

  addFolder(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const files: FileList | null = input.files;
    if (files && files.length > 1) {
      const data = this.dataSource.data
      for (let fileIndex = 0; fileIndex < files.length; ++fileIndex) {
        const file = files[fileIndex];
        console.log(file);
        const fileUpload: FileUpload = {
          fileName: file.name,
          // remove file name from path
          relativePath: file.webkitRelativePath.replace(new RegExp(file.name + '$'), ''),
          fileSize: file.size,
        }
        const fileDataIndex = data.length;
        data.push(fileUpload);
        this.uploadFile(file, fileDataIndex);
      }
      this.dataSource.data = data
    }
  }

  uploadFile(file: File, fileIndex: number): void {
    this.fileAnalysisService.analyzeFile(file).subscribe({
      error: (error: any) => {
        console.error(error);
      },
      next: (httpEvent: HttpEvent<FileAnalysis>) => {
        this.handleHttpEvent(httpEvent, fileIndex);
      }
    });
  }

  private handleHttpEvent(event: HttpEvent<FileAnalysis>, fileIndex: number): void {
    if (event.type === HttpEventType.UploadProgress) {
      if (event.total && event.total > 0.0) {
        this.dataSource.data[fileIndex].uploadProgress = Math.round(
          100 * (event.loaded / event.total)
        );
      }
    } else if(event.type === HttpEventType.Response) {
        if (event.body) {
          console.log(event.body);
          const fileData: FileInformation = {
            fileName: this.dataSource.data[fileIndex].fileName,
            relativePath: this.dataSource.data[fileIndex].relativePath,
            size: this.dataSource.data[fileIndex].fileSize,
            fileAnalysis: event.body,
          }
          this.fileAnalysisService.addFileInfo(fileData);
        }
    }
  }
}
