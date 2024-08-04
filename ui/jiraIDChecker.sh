HOOK_FILE=$1
COMMIT_MSG=`head -n1 $HOOK_FILE`
PATTERN="^[A-Z][A-Z0-9]+-[0-9]+[[:space:]][|]"
if [[ ! ${COMMIT_MSG} =~ $PATTERN ]]; then
  echo ""
  echo "    ERROR! Bad commit message. "
  echo "    '$COMMIT_MSG' is missing JIRA Ticket Number."
  echo "    example: 'TP-12345 | my commit'"
  echo ""
  exit 1
fi