import { Injectable, signal } from '@angular/core';
import { FileAnalysis, FileResult } from '../features/file-analysis/results';
import { formatFileSize } from '../shared/file-size.pipe';
import { FileUpload } from './file-analysis.service';

@Injectable({
  providedIn: 'root',
})
export class ResultsService {
  analysisDetails: { [key: string]: FileAnalysis } = {};
  fileResults = signal<FileResult[]>([]);

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
    this.fileResults.set([...this.fileResults(), fileResult]);
    this.analysisDetails[fileUpload.id] = analysis;
  }

  async get(id: string): Promise<FileAnalysis | undefined> {
    return this.analysisDetails[id];
  }
}
