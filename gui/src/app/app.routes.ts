import { Routes } from '@angular/router';
import { ResultsPageComponent } from './pages/results-page/results-page.component';
import { UploadPageComponent } from './pages/upload-page/upload-page.component';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'auswahl',
    pathMatch: 'full',
  },
  { path: 'auswahl', component: UploadPageComponent },
  { path: 'auswertung', component: ResultsPageComponent },
];
