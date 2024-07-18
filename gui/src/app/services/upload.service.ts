import { HttpEvent, HttpEventType } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { FileAnalysisService, FileUpload, ToolResults } from '../features/file-analysis/file-analysis.service';
import { NotificationService } from './notification.service';

@Injectable({
  providedIn: 'root',
})
export class UploadService {
  uploadInProgress = false;

  constructor(
    private router: Router,
    private notificationService: NotificationService,
    private fileAnalysis: FileAnalysisService,
  ) {
    this.fileAnalysis.getFileUploads().subscribe((fileUploads: FileUpload[]) => {
      if (fileUploads.length === 0 && this.uploadInProgress) {
        this.uploadInProgress = false;
        this.router.navigate(['auswertung']);
      }
    });
  }

  uploadFile(file: File, fileUpload: FileUpload): void {
    this.uploadInProgress = true;
    this.fileAnalysis.analyzeFile(file).subscribe({
      error: (error) => {
        fileUpload.error = error.statusText;
      },
      next: (httpEvent: HttpEvent<ToolResults>) => {
        this.handleHttpEvent(httpEvent, fileUpload);
      },
    });
  }

  private handleHttpEvent(event: HttpEvent<ToolResults>, fileUpload: FileUpload): void {
    if (event.type === HttpEventType.UploadProgress) {
      if (event.total && event.total > 0.0) {
        fileUpload.uploadProgress = Math.round(100 * (event.loaded / event.total));
      }
    } else if (event.type === HttpEventType.Response) {
      if (event.body) {
        this.fileAnalysis.addFileResult(fileUpload, event.body);
        this.notificationService.show('Formaterkennung und -validierung abgeschlossen: ' + fileUpload.fileName);
      }
    }
  }
}
