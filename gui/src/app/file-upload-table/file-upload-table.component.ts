// angular
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { HttpEventType, HttpEvent } from '@angular/common/http';
import { Router } from '@angular/router';

// material
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';

// project
import { 
  ToolResults, 
  FileUpload,
  FileAnalysisService,
} from '../file-analysis/file-analysis.service';
import { NotificationService } from 'src/app/utility/notification/notification.service';

@Component({
  selector: 'app-file-upload-table',
  templateUrl: './file-upload-table.component.html',
  styleUrls: ['./file-upload-table.component.scss'],
})
export class FileUploadTableComponent implements AfterViewInit {
  dataSource: MatTableDataSource<FileUpload>;
  displayedColumns: string[];
  uploadInProgress: boolean;

  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  constructor(
    private fileAnalysisService: FileAnalysisService,
    private notificationService: NotificationService,
    private router: Router,
  ) {
    this.uploadInProgress = false;
    this.dataSource = new MatTableDataSource<FileUpload>();
    this.displayedColumns = [
      'relativePath',
      'fileName',
      'fileSize',
      'uploadProgress',
      'verificationProgress',
    ];
    this.fileAnalysisService.getFileUploads().subscribe({
      // error can't occur --> no error handling
      next: (fileUploads: FileUpload[]) => {
        this.dataSource.data = fileUploads;
        if (fileUploads.length === 0 && this.uploadInProgress) {
          this.uploadInProgress = false;
          this.router.navigate(['auswertung']);
        }
      },
    });
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
      const fileUpload = this.fileAnalysisService.addFileUpload(
        file.name, 
        'Einzeldatei', 
        file.size,
      );
      this.uploadFile(file, fileUpload);
    }
  }

  addFolder(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const files: FileList | null = input.files;
    if (files && files.length > 1) {
      const data = this.dataSource.data
      for (let fileIndex = 0; fileIndex < files.length; ++fileIndex) {
        const file = files[fileIndex];
        const fileUpload = this.fileAnalysisService.addFileUpload(
          file.name, 
          file.webkitRelativePath.replace(new RegExp(file.name + '$'), ''), 
          file.size,
        );
        this.uploadFile(file, fileUpload);
      }
    }
  }

  uploadFile(file: File, fileUpload: FileUpload): void {
    this.uploadInProgress = true;
    this.fileAnalysisService.analyzeFile(file).subscribe({
      error: (error: any) => {
        console.error(error);
      },
      next: (httpEvent: HttpEvent<ToolResults>) => {
        this.handleHttpEvent(httpEvent, fileUpload);
      }
    });
  }

  private handleHttpEvent(event: HttpEvent<ToolResults>, fileUpload: FileUpload): void {
    if (event.type === HttpEventType.UploadProgress) {
      if (event.total && event.total > 0.0) {
        fileUpload.uploadProgress = Math.round(
          100 * (event.loaded / event.total)
        );
      }
    } else if(event.type === HttpEventType.Response) {
        if (event.body) {
          this.fileAnalysisService.addFileResult(
            fileUpload,
            event.body,
          );
          this.notificationService.show('Formaterkennung und -validierung abgeschlossen: ' 
            + fileUpload.fileName);
        }
    }
  }
}
