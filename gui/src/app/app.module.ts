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

import { DecimalPipe, registerLocaleData } from '@angular/common';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import localeDe from '@angular/common/locales/de';
import { LOCALE_ID, NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatPaginatorIntl, MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { FileAnalysisTableComponent } from './file-analysis/file-analysis-table/file-analysis-table.component';
import { FileOverviewComponent } from './file-overview/file-overview.component';
import { FileUploadTableComponent } from './file-upload-table/file-upload-table.component';
import { MainNavigationComponent } from './main-navigation/main-navigation.component';
import { ToolOutputComponent } from './tool-output/tool-output.component';
import { FileSizePipe } from './utility/formatting/file-size.pipe';
import { PrettyPrintCsvPipe } from './utility/formatting/pretty-print-csv.pipe';
import { PrettyPrintJsonPipe } from './utility/formatting/pretty-print-json.pipe';
import { FileFeaturePipe } from './utility/localization/file-attribut-de.pipe';
import { PaginatorDeService } from './utility/localization/paginator-de.service';

registerLocaleData(localeDe);

@NgModule({
  declarations: [
    AppComponent,
    FileAnalysisTableComponent,
    FileFeaturePipe,
    FileOverviewComponent,
    FileSizePipe,
    FileUploadTableComponent,
    MainNavigationComponent,
    PrettyPrintCsvPipe,
    PrettyPrintJsonPipe,
    ToolOutputComponent,
  ],
  bootstrap: [AppComponent],
  imports: [
    AppRoutingModule,
    BrowserAnimationsModule,
    BrowserModule,
    MatButtonModule,
    MatDialogModule,
    MatExpansionModule,
    MatIconModule,
    MatMenuModule,
    MatPaginatorModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatSortModule,
    MatTableModule,
    MatToolbarModule,
    MatTooltipModule,
  ],
  providers: [
    { provide: LOCALE_ID, useValue: 'de' },
    { provide: MatPaginatorIntl, useClass: PaginatorDeService },
    DecimalPipe,
    FileFeaturePipe,
    FileSizePipe,
    provideHttpClient(withInterceptorsFromDi()),
  ],
})
export class AppModule {}
