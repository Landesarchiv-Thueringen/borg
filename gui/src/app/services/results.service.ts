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
    const path = fileUpload.path;
    const fileSizeString = formatFileSize(fileUpload.fileSize);
    const fileResult: FileResult = {
      id: fileUpload.id,
      filename: fileUpload.filename,
      info: {
        path: { value: path },
        fileSize: {
          value: fileUpload.fileSize,
          displayString: fileSizeString,
        },
      },
      summary: analysis.summary,
    };
    this.fileResultsSubject.next([...this.fileResultsSubject.value, fileResult]);
    this.analysisDetails[fileUpload.id] = this.addBrowserInfo(analysis, path, fileSizeString);
  }

  addBrowserInfo(analysis: FileAnalysis, path: string, fileSizeString: string) {
    for (let set of analysis.featureSets) {
      set.features['general:path'] = {
        value: path,
        label: 'Pfad',
        supportingTools: ['browser'],
      };
      set.features['general:fileSize'] = {
        value: fileSizeString,
        label: 'Dateigröße',
        supportingTools: ['browser'],
      };
    }
    return analysis;
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
