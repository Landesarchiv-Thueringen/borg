import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { FileAnalysis, FileResult } from '../features/file-analysis/results';
import { formatFileSize } from '../shared/file-size.pipe';
import { FileUpload } from './file-analysis.service';

@Injectable({
  providedIn: 'root',
})
export class ResultsService {
  analysisDetails: { [key: string]: FileAnalysis } = {};
  fileResults: FileResult[] = [];
  fileResultsSubject = new BehaviorSubject<FileResult[]>(this.fileResults);

  add(fileUpload: FileUpload, analysis: FileAnalysis): void {
    const fileResult: FileResult = {
      id: fileUpload.id,
      filename: fileUpload.filename,
      info: {
        path: { value: fileUpload.path },
        fileSize: {
          value: fileUpload.fileSize,
          displayString: formatFileSize(fileUpload.fileSize),
        },
      },
      summary: analysis.summary,
    };
    this.fileResults.push(fileResult);
    this.fileResultsSubject.next(this.fileResults);
    this.analysisDetails[fileUpload.id] = analysis;
  }

  async get(id: string): Promise<FileAnalysis | undefined> {
    return this.analysisDetails[id];
  }

  getAll(): Observable<FileResult[]> {
    return this.fileResultsSubject.asObservable();
  }

  clear(): void {
    this.fileResults = [];
    this.fileResultsSubject.next(this.fileResults);
  }
}
