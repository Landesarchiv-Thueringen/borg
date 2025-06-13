import { HttpEvent, HttpEventType } from '@angular/common/http';
import { Injectable, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
import { v4 as uuid } from 'uuid';
import { FileAnalysis } from '../features/file-analysis/results';
import { FileAnalysisService, FileUpload } from './file-analysis.service';
import { NotificationService } from './notification.service';
import { ResultsService } from './results.service';

@Injectable({
  providedIn: 'root',
})
export class UploadService {
  private router = inject(Router);
  private notificationService = inject(NotificationService);
  private fileAnalysis = inject(FileAnalysisService);
  private results = inject(ResultsService);

  uploadInProgress = false;
  fileUploads = signal<FileUpload[]>([]);

  upload(file: File, fileUpload: FileUpload): void {
    this.uploadInProgress = true;
    this.fileAnalysis.analyzeFile(file).subscribe({
      error: (error) => {
        fileUpload.error = error.statusText;
      },
      next: (httpEvent: HttpEvent<FileAnalysis>) => {
        this.handleHttpEvent(httpEvent, fileUpload);
      },
    });
  }

  add(filename: string, path: string, fileSize: number): FileUpload {
    const fileUpload: FileUpload = {
      id: uuid(),
      filename: filename,
      path: path,
      fileSize: fileSize,
    };
    this.fileUploads.set([...this.fileUploads(), fileUpload]);
    return fileUpload;
  }

  private remove(fileUpload: FileUpload): void {
    let uploads = this.fileUploads();
    uploads = uploads.filter((upload) => {
      return upload.id !== fileUpload.id;
    });
    this.fileUploads.set([...uploads]);
    if (uploads.length === 0 && this.uploadInProgress) {
      this.uploadInProgress = false;
      this.router.navigate(['auswertung']);
    }
  }

  private handleHttpEvent(event: HttpEvent<FileAnalysis>, fileUpload: FileUpload): void {
    if (event.type === HttpEventType.UploadProgress) {
      if (event.total && event.total > 0.0) {
        fileUpload.uploadProgress = Math.round(100 * (event.loaded / event.total));
      }
    } else if (event.type === HttpEventType.Response) {
      if (event.body) {
        this.results.add(fileUpload, event.body);
        this.remove(fileUpload);
        this.notificationService.show(
          'Formaterkennung und -validierung abgeschlossen: ' + fileUpload.filename,
        );
      }
    }
  }
}
