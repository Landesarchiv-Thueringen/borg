import { HttpClient, HttpEvent } from '@angular/common/http';
import { Injectable } from '@angular/core';
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
  constructor(private httpClient: HttpClient) {}

  analyzeFile(file: File): Observable<HttpEvent<FileAnalysis>> {
    const formData = new FormData();
    formData.append('file', file);
    return this.httpClient.post<FileAnalysis>('/analyze-file', formData, {
      reportProgress: true,
      observe: 'events',
    });
  }
}
