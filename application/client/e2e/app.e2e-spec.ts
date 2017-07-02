import { AlfPage } from './app.po';

describe('alf App', () => {
  let page: AlfPage;

  beforeEach(() => {
    page = new AlfPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!!');
  });
});
