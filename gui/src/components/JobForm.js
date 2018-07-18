import axios from "axios";
import React from "react";
import { withFormik, Form } from "formik";
import { withStyles } from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import { TextField } from "@material-ui/core";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

const styles = theme => ({
  textfield: {
    padding: "10px 0px",
    width: "200px"
  },
  jsonfield: {
    padding: "10px 0px",
    width: "700px"
  },

  card: {
    paddingBottom: theme.spacing.unit * 2
  }
});

const App = ({ values, isSubmitting, classes, handleChange }) => (
  <div style={{ width: "50%" }}>
    <br />
    <Form style={{ left: "50%", position: "relative" }}>
      <div>
        <TextField
          id="textarea"
          label="Paste JSON"
          onChange={handleChange}
          placeholder="Paste JSON"
          multiline
          className={classes.jsonfield}
          margin="normal"
        />
      </div>
      <Button color="primary" type="submit" disabled={isSubmitting}>
        Build Job
      </Button>
      <ToastContainer  />
    </Form>
  </div>
);

const JobForm = withFormik({
  mapPropsToValues({ json }) {
    return {
      json: json || ""
    };
  },
  handleSubmit(values) {
    axios
      .post("/v2/specs", JSON.parse(values.textarea), {
        headers: {
          "Content-Type": "application/json"
        },
        auth: {
          username: "chainlink",
          password: "twochains"
        }
      })
      .then(res =>
        toast.success(`Job ${res.data.id} created`, {
          position: toast.POSITION.BOTTOM_RIGHT
        })
      );
  }
})(App);

export default withStyles(styles)(JobForm);
