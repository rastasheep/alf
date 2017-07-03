import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { EditorModule } from './editor/editor.module';

import { AppComponent } from './app.component';

export const ROUTES: Routes = [];

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    RouterModule.forRoot(ROUTES),
    EditorModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {}
