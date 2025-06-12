import { HttpClient, HttpEvent } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { FileAnalysis } from '../features/file-analysis/results';

export interface FileUpload {
  id: string;
  filename: string;
  path: string;
  fileSize: number;
  uploadProgress?: number;
  error?: string;
}

@Injectable({
  providedIn: 'root',
})
export class FileAnalysisService {
  private httpClient = inject(HttpClient);

  analyzeFile(file: File): Observable<HttpEvent<FileAnalysis>> {
    const formData = new FormData();
    formData.append('file', file);
    return this.httpClient.post<FileAnalysis>('/api/analyze', formData, {
      reportProgress: true,
      observe: 'events',
    });
  }
}
