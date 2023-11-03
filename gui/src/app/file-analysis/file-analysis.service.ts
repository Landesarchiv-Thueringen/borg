// angular
import { Injectable } from '@angular/core';
import { HttpClient, HttpEvent } from '@angular/common/http';
import { Observable } from 'rxjs';

// project
import { environment } from '../../environments/environment';

// utility
import { BehaviorSubject } from 'rxjs';

export interface FileInformation {
  fileName: string;
  relativePath?: string;
  size: string;
  fileAnalysis: FileAnalysis;
}

export interface FileAnalysis {
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
  files: FileInformation[];
  fileInfos: BehaviorSubject<FileInformation[]>;
  featureOrder: Map<string, number>;

  constructor(private httpClient: HttpClient) {
    this.files = [];
    this.fileInfos = new BehaviorSubject<FileInformation[]>(this.files);
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

  analyzeFile(file: File): Observable<HttpEvent<FileAnalysis>> {
    const formData = new FormData();
    formData.append('file', file);
    console.log(environment.apiEndpoint);
    return this.httpClient.post<FileAnalysis>(environment.apiEndpoint, formData, {
      reportProgress: true,
      observe: 'events'
    });
  }

  addFileInfo(i: FileInformation): void {
    this.files.push(i);
    this.fileInfos.next(this.files);
  }

  getFileInfo(): Observable<FileInformation[]> {
    return this.fileInfos.asObservable();
  }

  getFeatureOrder(): Map<string, number> {
    return this.featureOrder;
  }
}
