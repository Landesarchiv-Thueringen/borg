import { HttpEvent, HttpEventType } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject, Observable } from 'rxjs';
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
  fileUploads: FileUpload[] = [];
  fileUploadsSubject = new BehaviorSubject<FileUpload[]>(this.fileUploads);

  constructor() {
    this.getAll().subscribe((fileUploads: FileUpload[]) => {
      if (fileUploads.length === 0 && this.uploadInProgress) {
        this.uploadInProgress = false;
        this.router.navigate(['auswertung']);
      }
    });
  }

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
    this.fileUploads.push(fileUpload);
    this.fileUploadsSubject.next(this.fileUploads);
    return fileUpload;
  }

  getAll(): Observable<FileUpload[]> {
    return this.fileUploadsSubject.asObservable();
  }

  private remove(fileUpload: FileUpload): void {
    this.fileUploads = this.fileUploads.filter((upload: FileUpload) => {
      return upload.id !== fileUpload.id;
    });
    this.fileUploadsSubject.next(this.fileUploads);
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
