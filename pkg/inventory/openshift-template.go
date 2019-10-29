package inventory

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type OpenShiftTemplate struct {
	Template   string            `yaml:"template"`
	Params     map[string]string `yaml:"params"`
	ParamFiles []string          `yaml:"paramFiles"`
}

func (ot *OpenShiftTemplate) Process(ns *string, r *Resource) error {

	p := r.Prefix + "/" + ot.Template
	abs, err := filepath.Abs(p)
	if err != nil {
		return err
	}

	// check which processing mode to use
	tp, err := os.Stat(abs)
	if err != nil {
		return err
	}
	val, ok := ot.Params["param_dir"]
	if tp.IsDir() {

		// get all template files in diectory
		var templates []string
		err := filepath.Walk(abs, func(path string, info os.FileInfo, err error) error {
			templates = append(templates, path)
			return nil
		})
		if err != nil {
			return err
		}

		if ok {
			// Case 1: User has passed a directory of templates, and a directory of parameters.
			// We will expect a one to one mapping of template file to parameter file.
			// get param file of the same name
			for _, template := range templates {
				// process template and file
				err = processOneTemplate(template, []string{val + "/" + filepath.Base(template)}, make(map[string]string), r.Output+"/"+string(r.Action))
				if err != nil {
					return err
				}
			}
		} else {
			// Case 2: User has passed a directory of templates, and a single set of params
			for _, template := range templates {
				// process template and file
				err = processOneTemplate(template, ot.ParamFiles, ot.Params, r.Output+"/"+string(r.Action))
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	if ok {
		// Case 3: User has passed a directory of params and a single template. We will
		// process the template once for each param file
		// get all template files in diectory
		var params []string
		err := filepath.Walk(val, func(path string, info os.FileInfo, err error) error {
			params = append(params, path)
			return nil
		})
		if err != nil {
			return err
		}
		for _, param := range params {
			// process template and file
			err = processOneTemplate(abs, []string{param}, ot.Params, r.Output+"/"+string(r.Action))
			if err != nil {
				return err
			}
		}
		return nil
	}
	// Case 4: One template, one set of params
	err = processOneTemplate(abs, ot.ParamFiles, ot.Params, r.Output+"/"+string(r.Action))
	if err != nil {
		return err
	}

	return nil
}

func processOneTemplate(tpl string, pF []string, ps map[string]string, r *Resource) error {
	// oc process -f template-file -p PARAM=foo --param-file
	cmdArgs := []string{"process", "--local", "-f", tpl}
	for key, param := range ps {
		cmdArgs = append(cmdArgs, "-p", key+"="+param)
	}
	for _, pf := range pF {
		cmdArgs = append(cmdArgs, "--param-file", pf)
	}
	cmd := exec.Command("oc", cmdArgs...)
	log.Printf("Running command: %s\n", cmd.Args)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s\n", stdoutStderr)
		return err
	}

	outDir := r.Prefix + "/" + r.Action
	// write resulting resource to file
	output_dir := filepath.Clean(outDir)
	out, err := os.Create(output_dir + "/" + filepath.Base(tpl))
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	log.Printf("wrote %s\n", out.Name())
	_, err = out.Write(stdoutStderr)
	if err != nil {
		return err
	}

	return nil

}
