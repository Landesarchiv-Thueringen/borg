/* BorgFormat - File format identification and validation
 * Copyright (C) 2024 Landesarchiv Thüringen
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import { HttpClient, HttpEvent } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { v4 as uuid } from 'uuid';

export interface FileUpload {
  id: string;
  fileName: string;
  relativePath: string;
  fileSize: number;
  uploadProgress?: number;
}

export interface FileResult {
  id: string;
  fileName: string;
  relativePath?: string;
  fileSize: number;
  toolResults: ToolResults;
}

export interface ToolResults {
  fileIdentificationResults: ToolResult[] | null;
  fileValidationResults: ToolResult[] | null;
  summary: Summary;
}

export interface ToolResult {
  toolName: string;
  toolVersion: string;
  toolOutput: string;
  outputFormat: 'text' | 'json' | 'csv';
  extractedFeatures: { [key: string]: string };
  error: string | null;
}

const OVERVIEW_FEATURES = [
  'relativePath',
  'fileName',
  'fileSize',
  'puid',
  'mimeType',
  'formatVersion',
  'valid',
] as const;
export type OverviewFeature = (typeof OVERVIEW_FEATURES)[number];

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
  providedIn: 'root',
})
export class FileAnalysisService {
  fileUploads: FileUpload[] = [];
  fileResults: FileResult[] = [];
  fileUploadsSubject = new BehaviorSubject<FileUpload[]>(this.fileUploads);
  fileResultsSubject = new BehaviorSubject<FileResult[]>(this.fileResults);
  readonly featureOrder = new Map<string, number>([
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

  constructor(private httpClient: HttpClient) {}

  analyzeFile(file: File): Observable<HttpEvent<ToolResults>> {
    const formData = new FormData();
    formData.append('file', file);
    return this.httpClient.post<ToolResults>("/analyze-file", formData, {
      reportProgress: true,
      observe: 'events',
    });
  }

  addFileResult(fileUpload: FileUpload, toolResults: ToolResults): void {
    const fileResult: FileResult = {
      id: fileUpload.id,
      fileName: fileUpload.fileName,
      relativePath: fileUpload.relativePath,
      fileSize: fileUpload.fileSize,
      toolResults: toolResults,
    };
    this.fileResults.push(fileResult);
    this.fileResultsSubject.next(this.fileResults);
    this.removeFileUpload(fileUpload);
  }

  getFileResult(id: string): FileResult | undefined {
    return this.fileResults.find((fileResult) => fileResult.id === id);
  }

  getFileResults(): Observable<FileResult[]> {
    return this.fileResultsSubject.asObservable();
  }

  clearFileResults(): void {
    this.fileResults = [];
    this.fileResultsSubject.next(this.fileResults);
  }

  addFileUpload(fileName: string, relativePath: string, fileSize: number): FileUpload {
    const fileUpload: FileUpload = {
      id: uuid(),
      fileName: fileName,
      relativePath: relativePath,
      fileSize: fileSize,
    };
    this.fileUploads.push(fileUpload);
    this.fileUploadsSubject.next(this.fileUploads);
    return fileUpload;
  }

  getFileUploads(): Observable<FileUpload[]> {
    return this.fileUploadsSubject.asObservable();
  }

  removeFileUpload(fileUpload: FileUpload): void {
    this.fileUploads = this.fileUploads.filter((upload: FileUpload) => {
      return upload.id !== fileUpload.id;
    });
    this.fileUploadsSubject.next(this.fileUploads);
  }

  getFeatureOrder(): Map<string, number> {
    return this.featureOrder;
  }

  /** Sorts feature keys and removes duplicates. */
  sortFeatures(features: string[]): string[] {
    features = [...new Set(features)];
    return features.sort((f1: string, f2: string) => {
      const featureOrder = this.getFeatureOrder();
      let orderF1: number | undefined = featureOrder.get(f1);
      if (!orderF1) {
        orderF1 = featureOrder.get('');
      }
      let orderF2: number | undefined = featureOrder.get(f2);
      if (!orderF2) {
        orderF2 = featureOrder.get('');
      }
      if (orderF1! < orderF2!) {
        return -1;
      } else if (orderF1! > orderF2!) {
        return 1;
      }
      return 0;
    });
  }

  isOverviewFeature(feature: string): feature is OverviewFeature {
    return (OVERVIEW_FEATURES as readonly string[]).includes(feature);
  }
}
