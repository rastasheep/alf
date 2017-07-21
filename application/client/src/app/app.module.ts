import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardModule } from './dashboard/dashboard.module';

import { AppComponent } from './app.component';

export const ROUTES: Routes = [
  { path: '', redirectTo: 'editor', pathMatch: 'full' },
];

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    RouterModule.forRoot(
      ROUTES,
      { enableTracing: true } // debugging purposes only
   ),
    DashboardModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {}
