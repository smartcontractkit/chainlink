import axios from "axios";
import React from "react";
import { withFormik, Form } from "formik";
import * as Yup from "yup";
import { withStyles } from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import { TextField, Typography } from "@material-ui/core";

const styles = theme => ({
  textfield: {
    padding: "10px 0px",
    width: "270px"
  },
  card: {
    paddingBottom: theme.spacing.unit * 2
  }
});

const App = ({ values, errors, touched, isSubmitting, classes, handleChange }) => (
  <div style={{}}>
    <br />
    <Form style={{ position: "relative", textAlign: "center" }} noValidate>
      <div>
        {touched.name && errors.name && <Typography color="error">{errors.name}</Typography>}
        <TextField
          onChange={handleChange}
          className={classes.textfield}
          label="Type Bridge Name"
          type="name"
          name="name"
          placeholder="name"
        />
      </div>
      <div>
        {touched.url && errors.url && <Typography color="error">{errors.url}</Typography>}
        <TextField
          onChange={handleChange}
          className={classes.textfield}
          label="Type Bridge URL"
          type="url"
          name="url"
          placeholder="url"
        />
      </div>
      <div>
        {touched.confirmations && errors.confirmations && <Typography color="error">{errors.confirmations}</Typography>}
        <TextField
          onChange={handleChange}
          className={classes.textfield}
          label="Type Confirmations"
          type="confirmations"
          name="confirmations"
          placeholder="confirmations"
        />
      </div>
      <Button color="primary" type="submit" disabled={isSubmitting}>
        Build Bridge
      </Button>
    </Form>
  </div>
);

const BridgeForm = withFormik({
  mapPropsToValues({ name, url, confirmations }) {
    return {
      name: name || "",
      url: url || "",
      confirmations: confirmations || ""
    };
  },
  validationSchema: Yup.object().shape({
    name: Yup.string().required("Name is required"),
    url: Yup.string()
      .url("Should be a valid link")
      .required("URL is required"),
    confirmations: Yup.number()
      .positive("Should be a positive number")
      .typeError("Should be a number")
  }),
  handleSubmit(values) {
    const formattedValues = JSON.parse(JSON.stringify(values).replace("confirmations", "defaultConfirmations"));
    formattedValues.defaultConfirmations = parseInt(formattedValues.defaultConfirmations) || 0
    axios
      .post("/v2/bridge_types", formattedValues, {
        headers: {
          "Content-Type": "application/json"
        },
        auth: {
          username: "chainlink",
          password: "twochains"
        }
      })
      .then(res => console.log(res));
  }
})(App);

export default withStyles(styles)(BridgeForm);
