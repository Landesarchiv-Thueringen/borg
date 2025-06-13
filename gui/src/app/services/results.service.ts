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
  get fileResults() {
    return this.fileResultsSubject.value;
  }
  private fileResultsSubject = new BehaviorSubject<FileResult[]>([]);

  add(fileUpload: FileUpload, analysis: FileAnalysis): void {
    const fileResult: FileResult = {
      id: fileUpload.id,
      filename: fileUpload.filename,
      summary: analysis.summary,
      additionalMetadata: {
        'general:path': {
          value: fileUpload.path,
          label: 'Pfad',
          supportingTools: ['browser'],
        },
        'general:fileSize': {
          value: formatFileSize(fileUpload.fileSize),
          label: 'Dateigröße',
          supportingTools: ['browser'],
        },
      },
    };
    this.fileResultsSubject.next([...this.fileResultsSubject.value, fileResult]);
    this.analysisDetails[fileUpload.id] = analysis;
  }

  async get(id: string): Promise<FileAnalysis | undefined> {
    return this.analysisDetails[id];
  }

  getAll(): Observable<FileResult[]> {
    return this.fileResultsSubject.asObservable();
  }

  clear(): void {
    this.fileResultsSubject.next([]);
  }
}
