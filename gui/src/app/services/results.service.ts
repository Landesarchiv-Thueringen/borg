import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { FileResult, ToolResults } from '../features/file-analysis/results';
import { FileUpload } from './file-analysis.service';

@Injectable({
  providedIn: 'root',
})
export class ResultsService {
  fileResults: FileResult[] = [];
  fileResultsSubject = new BehaviorSubject<FileResult[]>(this.fileResults);

  add(fileUpload: FileUpload, toolResults: ToolResults): void {
    const fileResult: FileResult = {
      id: fileUpload.id,
      fileName: fileUpload.fileName,
      relativePath: fileUpload.relativePath,
      fileSize: fileUpload.fileSize,
      toolResults: toolResults,
    };
    this.fileResults.push(fileResult);
    this.fileResultsSubject.next(this.fileResults);
  }

  async get(id: string): Promise<FileResult | undefined> {
    return this.fileResults.find((fileResult) => fileResult.id === id);
  }

  getAll(): Observable<FileResult[]> {
    return this.fileResultsSubject.asObservable();
  }

  clear(): void {
    this.fileResults = [];
    this.fileResultsSubject.next(this.fileResults);
  }
}
