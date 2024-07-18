import { Routes } from '@angular/router';
import { FileAnalysisTableComponent } from './file-analysis/file-analysis-table/file-analysis-table.component';
import { FileUploadTableComponent } from './file-upload-table/file-upload-table.component';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'auswahl',
    pathMatch: 'full',
  },
  { path: 'auswahl', component: FileUploadTableComponent },
  { path: 'auswertung', component: FileAnalysisTableComponent },
];
