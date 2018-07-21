/* eslint-env jest */
import React from "react";
import { shallow } from "enzyme";
import { Conn as Create } from "containers/Create";

const classes = {};
const mountCreatePage = props => shallow(<Create classes={classes} {...props} />);

describe("containers/Create", () => {
  it("renders the create page focused on bridge creation", () => {
    expect.assertions(3);
    const props = { location: { state: { tab: 0 } } };
    const wrapper = mountCreatePage(props).render();
    expect(wrapper.text()).toContain("Create Bridge");
    expect(wrapper.text()).toContain("Build Bridge");
    expect(wrapper.text()).toContain("Type Confirmations");
  });

  it("renders the create page focused on job creation", () => {
    expect.assertions(3);
    const props = { location: { state: { tab: 1 } } };
    let wrapper = mountCreatePage(props);
    wrapper.instance().forceUpdate();
    wrapper.update();
    wrapper = wrapper.render();
    expect(wrapper.text()).toContain("Create Job");
    expect(wrapper.text()).toContain("Build Job");
    expect(wrapper.text()).toContain("Paste JSON");
  });
});
