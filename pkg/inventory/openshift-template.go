package inventory

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type OpenShiftTemplate struct {
	Template   string            `yaml:"template"`
	Params     map[string]string `yaml:"params"`
	ParamFiles []string          `yaml:"paramFiles"`
	ParamDir   string            `yaml:"paramDir"`
}

func (ot *OpenShiftTemplate) Process(ns *string, r *Resource) error {

	// Pre-process to convert to absolute pathing
	abs, err := filepath.Abs(filepath.Join(r.Prefix, ot.Template))
	if err != nil {
		return err
	}
	ot.Template = abs

	ok := ot.ParamDir != ""
	if ok {
		abs, err = filepath.Abs(filepath.Join(r.Prefix, ot.ParamDir))
		if err != nil {
			return err
		}
		ot.ParamDir = abs
	}

	var pFs []string
	for _, pF := range ot.ParamFiles {
		abs, err = filepath.Abs(filepath.Join(r.Prefix, pF))
		if err != nil {
			return err
		}
		pFs = append(pFs, abs)
	}
	ot.ParamFiles = pFs

	// check which processing mode to use
	tp, err := os.Stat(ot.Template)
	if err != nil {
		return err
	}

	if tp.IsDir() {

		// get all template files in diectory
		var templates []string
		err := filepath.Walk(ot.Template, appendFile(&templates))
		if err != nil {
			return err
		}

		if ok {
			// Case 1: User has passed a directory of templates, and a directory of parameters.
			// We will expect a one to one mapping of template file to parameter file.
			// get param file of the same name
			fmt.Printf("Found template directory %s and param directory %s\n", ot.Template, ot.ParamDir)
			for _, template := range templates {
				// process template and file
				ext := filepath.Ext(template)
				filename := filepath.Base(template)
				newVal := filepath.Join(ot.ParamDir, strings.Replace(filename, ext, "", -1))
				err = processOneTemplate(template, []string{newVal}, ot.Params, r, "")
				if err != nil {
					return err
				}
			}
		} else {
			// Case 2: User has passed a directory of templates, and a single set of params
			fmt.Printf("Found template directory %s and one set of params\n", ot.Template)
			for _, template := range templates {
				// process template and file
				err = processOneTemplate(template, ot.ParamFiles, ot.Params, r, "")
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
		fmt.Printf("Found template %s and a param directory %s\n", ot.Template, ot.ParamDir)
		var params []string
		err := filepath.Walk(ot.ParamDir, appendFile(&params))
		if err != nil {
			return err
		}
		for _, param := range params {
			// process template and file
			err = processOneTemplate(ot.Template, []string{param}, ot.Params, r, filepath.Base(param)+".yaml")
			if err != nil {
				return err
			}
		}
		return nil
	}
	// Case 4: One template, one set of params
	fmt.Printf("Found template %s and one set of params\n", ot.Template)
	err = processOneTemplate(ot.Template, ot.ParamFiles, ot.Params, r, "")
	if err != nil {
		return err
	}

	return nil
}

func processOneTemplate(tpl string, pF []string, ps map[string]string, r *Resource, outName string) error {
	if outName == "" {
		outName = filepath.Base(tpl)
	}
	// oc process -f template-file -p PARAM=foo --param-file
	cmdArgs := []string{"process", "--local", "--ignore-unknown-parameters", "-f", tpl}
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

	// write resulting resource to file
	outputDir := filepath.Join(r.Output, r.Action)
	out, err := os.Create(filepath.Join(outputDir, outName))
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

func appendFile(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		new, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !new.IsDir() {
			*files = append(*files, path)
		}
		return nil
	}
}
