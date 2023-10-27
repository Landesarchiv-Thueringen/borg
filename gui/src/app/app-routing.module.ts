// angular
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

// project
import { FileUploadTableComponent } from './file-upload-table/file-upload-table.component';
import { MainNavigationComponent } from './main-navigation/main-navigation.component';

const routes: Routes = [
  {
    path: '',  component: MainNavigationComponent,
    children: [
      { path: "upload", component: FileUploadTableComponent }
    ],
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
