// angular
import { Injectable } from '@angular/core';
import { HttpClient, HttpEvent } from '@angular/common/http';
import { Observable } from 'rxjs';

// project
import { environment } from '../../environments/environment';

// utility
import { BehaviorSubject } from 'rxjs';
import { v4 as uuidv4 } from 'uuid';

export interface FileUpload {
  fileName: string;
  relativePath: string;
  fileSize: number;
  uploadProgress?: number;
}

export interface FileResult {
  fileName: string;
  relativePath?: string;
  fileSize: number;
  toolResults: ToolResults;
}

export interface ToolResults {
  fileIdentificationResults: ToolResult[];
  fileValidationResults: ToolResult[];
  summary: Summary;
}

export interface ToolResult {
  toolName: string;
  toolVersion: string;
  toolOutput: string;
}

export interface Summary {
  [key: string]: Feature;
}

export interface Feature {
  key: string;
  values: FeatureValue[];
}

export interface FeatureValue {
  value: string;
  score: number;
  tools: ToolConfidence[];
}

export interface ToolConfidence {
  confidence: number;
  toolName: string;
}

@Injectable({
  providedIn: 'root'
})
export class FileAnalysisService {
  fileUploads: FileUpload[];
  fileResults: FileResult[];
  fileUploadsSubject: BehaviorSubject<FileUpload[]>;
  fileResultsSubject: BehaviorSubject<FileResult[]>;
  featureOrder: Map<string, number>;

  constructor(private httpClient: HttpClient) {
    this.fileUploads = [];
    this.fileResults = [];
    this.fileUploadsSubject = new BehaviorSubject<FileUpload[]>(this.fileUploads);
    this.fileResultsSubject = new BehaviorSubject<FileResult[]>(this.fileResults);
    this.featureOrder = new Map<string, number>([
      ['relativePath', 1],
      ['fileName', 2],
      ['fileSize', 3],
      ['puid', 4],
      ['mimeType', 5],
      ['formatVersion', 6],
      ['encoding', 7],
      ['', 101],
      ['wellFormed', 1001],
      ['valid', 1002],
    ]);
  }

  analyzeFile(file: File): Observable<HttpEvent<ToolResults>> {
    const formData = new FormData();
    formData.append('file', file);
    return this.httpClient.post<ToolResults>(environment.apiEndpoint, formData, {
      reportProgress: true,
      observe: 'events'
    });
  }

  addFileResult(
    fileName: string, 
    relativePath: string, 
    fileSize: number,
    toolResults: ToolResults,
  ): void {
    const fileResult: FileResult = {
      fileName: fileName,
      relativePath: relativePath,
      fileSize: fileSize,
      toolResults: toolResults,
    }
    this.fileResults.push(fileResult);
    this.fileResultsSubject.next(this.fileResults);
  }

  getFileResults(): Observable<FileResult[]> {
    return this.fileResultsSubject.asObservable();
  }

  addFileUpload(fileName: string, relativePath: string, fileSize: number): FileUpload {
    const fileUpload: FileUpload = {
      fileName: fileName,
      relativePath: relativePath,
      fileSize: fileSize,
    }
    this.fileUploads.push(fileUpload)
    this.fileUploadsSubject.next(this.fileUploads);
    return fileUpload;
  }

  getFileUploads(): Observable<FileUpload[]> {
    return this.fileUploadsSubject.asObservable();
  }

  getFeatureOrder(): Map<string, number> {
    return this.featureOrder;
  }
}
